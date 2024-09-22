package controller 

import (
    "net/http"
    "log"

    "github.com/ty/kshared/conf"
    "github.com/ty/kshared/model"
    "github.com/ty/kshared/controller/utils"
)

func registerRecvNewMdHandle() {
    http.HandleFunc("/upload", RecvNewMdFile)
}

func RecvNewMdFile(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    if err := r.ParseMultipartForm(conf.MaxRecvFileMem()); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // read uid from cookie
    cookie, err := r.Cookie("uid");
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    } 
    
    a, err := model.NewArticleByItem(
        r.PostFormValue("atitle"),
        r.PostFormValue("atype"),
        r.PostFormValue("alabel"),
        cookie.Value,
    )
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // recv file 
    file, _, err := r.FormFile("uploadfile")
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    defer file.Close()
    t, err := utils.NewTemp()
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    defer t.Close()
    if err := t.Copy(file); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if err := a.Insert(); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    t.Save(a.Id)

    // response home
    handleHome(w, r) // Oh! Cao
}

