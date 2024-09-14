package controller 

import (
    "net/http"
    "log"

    "github.com/ty/kshared/model"
    "github.com/ty/kshared/view"
)

func registerArchiveHandle() {
    http.HandleFunc("/archive", getArchive)
}



func getArchive(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    infos, err := model.GetTotalArticle()
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    err = view.ExecuteTmpl(w, &view.TmplArgs{
        Type: "archive",
        Value: infos,
    })
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
}
