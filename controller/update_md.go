package controller 

import (
    "net/http"
    "log"

    "github.com/ty/kshared/model"
    "github.com/ty/kshared/conf"
    "github.com/ty/kshared/controller/utils"
)

func registerUpdateMdHandle() {
    http.HandleFunc("/update", updateMdFile)
}


// this handle is only update file content and update_time
func updateMdFile(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query()["id"]
    if r.Method != http.MethodPost || len(id) != 1 {
        log.Println("Bad method or len != 1")
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

    // get article 
    a, err := model.NewArticleById(id[0], cookie.Value) 
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // recv file
    file, _, err := r.FormFile("updatefile")
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

    if err := a.Update(); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    t.Save(a.Id)
    log.Println("updatefile ok:", a.UpdateTime)
    w.WriteHeader(http.StatusOK)
}
