package main 

import (
    "net/http"

    "github.com/ty/kshared/controller"
    "github.com/ty/kshared/middleware"
)

func main() {
    controller.RegisterHandler()

    s := http.Server {
        Addr: "192.168.18.128:8888",
        Handler: &middleware.CrosMiddleawre{}, 
    }

    println("Server is running...")
    s.ListenAndServe()
}
