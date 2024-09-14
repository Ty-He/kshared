package controller 

import (
    "net/http"
    "log"
    
    "github.com/ty/kshared/model"
    "github.com/ty/kshared/view"
    "github.com/ty/kshared/conf"
)

func registerCategoryHandle() {
    http.HandleFunc("/category", getCategory)
}


func getCategory(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // get information about category from conf
    categoryContent, err := model.FilterArticleByCategory(conf.Category())
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    err = view.ExecuteTmpl(w, &view.TmplArgs{
        Type: "category",
        Value: categoryContent,
    })

    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
}
