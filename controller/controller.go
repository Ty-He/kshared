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

    // FileServer
    http.Handle("/css/", http.FileServer(http.Dir("resource/web")))
    http.Handle("/js/", http.FileServer(http.Dir("resource/web")))
    http.Handle("/img/", http.FileServer(http.Dir("resource/web")))

}
