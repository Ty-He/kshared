package controller 

import (
    "net/http"
    "log"

    "github.com/ty/kshared/model"
    "github.com/ty/kshared/view"
)

func registerSearchHandle() {
    http.HandleFunc("/search", handleSearch)
}


func handleSearch(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    search := r.URL.Query().Get("tag")
    
    t := model.Tag(search)

    items, err := t.GetArticle()
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    tagVal := &struct {
        TagName string
        Value any
    } {
        TagName: search,
        Value: items,
    }

    err = view.ExecuteTmpl(w, &view.TmplArgs{
        Type: "tag",
        Value: tagVal,
    })
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

}
