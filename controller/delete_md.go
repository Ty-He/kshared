package controller 

import (
    "net/http"
    "log"

    "github.com/ty/kshared/model"
)

func registerDeleteMdHandle() {
    http.HandleFunc("/delete", deleteMdFile)

}
func deleteMdFile(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query()["id"]
    if r.Method != http.MethodDelete || len(id) != 1 {
        log.Println("Bad method or len != 1")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    
    cookie, err := r.Cookie("uid");
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    } 

    a, err := model.NewArticleById(id[0], cookie.Value)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if err := a.Delete(); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // ok 
    w.WriteHeader(http.StatusOK)
}
