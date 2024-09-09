package controller 

import (
    "net/http"
    "os"
    "log"
)

func registerHomeHanle() {
    http.HandleFunc("/home", handleHome)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        content, err := os.ReadFile("resource/article/start.md")
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        w.Write(content)
    } else {
        w.WriteHeader(http.StatusBadRequest)
    }
}
