package main 

import (
    "fmt"
    "net/http"
    "log"
    "errors"
    "os"
    "os/signal"
    "syscall"
    "context"
    "time"

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

    go func() {
        fmt.Printf("Http server [%s] is running at local port: %d ...\n", conf.ServerName(), conf.Port())
        if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
            log.Fatalln(err)
        }
        fmt.Printf("Http server [%s] begin exiting.\n", conf.ServerName())
    }()

    signSet := make(chan os.Signal, 1)
    signal.Notify(signSet, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

    <-signSet
    fmt.Println(" Catch exiting signal.")
    
    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

    if err := s.Shutdown(ctx); err != nil {
        log.Println(err)
        return
    }
    fmt.Println("Server graceful exit.")
}

