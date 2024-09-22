package controller 

import (
    "net/http"
    "log"

    "github.com/ty/kshared/model"
    "github.com/ty/kshared/view"
)

func registerHomeHandle() {
    http.HandleFunc("/", handleHome)
}

// this func should allow all method
func handleHome(w http.ResponseWriter, r *http.Request) {
    // 1 get recent article list
    infos, err := model.GetLatestArticles()
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    // 2 load template 
    // 3 write reponse
    err = view.ExecuteTmpl(w, &view.TmplArgs{
        Type: "home",
        Value: infos,
    })
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
}
