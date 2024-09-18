package controller 

import "net/http"

func RegisterHandler() {
    registerHomeHandle()
    registerMdHandle()
    registerArchiveHandle()
    registerCategoryHandle()
    registerOnlineHandle()
    registerRecvNewMdHandle()
    registerUpdateMdHandle()
    registerDeleteMdHandle()
    registerCommentHandler()

    // FileServer
    http.Handle("/css/", http.FileServer(http.Dir("resource/web")))
    http.Handle("/js/", http.FileServer(http.Dir("resource/web")))
    http.Handle("/img/", http.FileServer(http.Dir("resource/web")))
    // for favicon
    http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "resource/web/favicon.ico")
    })
}
