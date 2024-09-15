package main 

import (
    "fmt"
    "net/http"
    "log"

    "github.com/ty/kshared/controller"
    "github.com/ty/kshared/middleware"
    "github.com/ty/kshared/conf"
)

func main() {
    log.SetFlags(log.Lshortfile)
    controller.RegisterHandler()

    s := http.Server {
        Addr: fmt.Sprintf("192.168.18.128:%d", conf.Port()),
        Handler: &middleware.CrosMiddleawre{}, 
    }


    fmt.Printf("Http server [%s] is running at local port: %d ...\n", conf.ServerName(), conf.Port())
    s.ListenAndServe()
}
