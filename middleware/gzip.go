package middleware 

import (
    "compress/gzip"
    "net/http"
    "strings"
)

type gzipResponseWriter struct {
    http.ResponseWriter 
    g *gzip.Writer
}

// override
func (w *gzipResponseWriter) Writer(p []byte) (int, error) {
    if w.g != nil {
        return w.g.Write(p)
    }
    // should not call self
    return w.ResponseWriter.Write(p)
}

type GzipMiddleware struct {
    next http.Handler
}

func (g *GzipMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
        w.Header().Set("Content-Encoding", "gzip")
        gw := gzip.NewWriter(w)
        defer gw.Close()

        nw := &gzipResponseWriter{
            ResponseWriter: w,
            g: gw,
        }
        g.next.ServeHTTP(nw, r)
    } else {
        g.next.ServeHTTP(w, r)
    }
}


