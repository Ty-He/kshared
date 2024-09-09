package middleware 

import (
    "log"
    "net/http"
)

type CrosMiddleawre struct {
    Next http.Handler
}

func (c *CrosMiddleawre) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if c.Next == nil {
        c.Next = http.DefaultServeMux
    }

    log.Println("handler request: ", r.Method, r.URL)

    // CROS
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // preflight request
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    c.Next.ServeHTTP(w, r)    
}

