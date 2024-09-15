package controller 

import (
    "net/http"
    "log"
    "os"
    "io"
    "fmt"

    "github.com/ty/kshared/conf"
    "github.com/ty/kshared/model"
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
    
    a, err := model.NewArticle(
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
    outfile, err := os.Create("resource/article/0.md")
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    defer outfile.Close()
    
    if _, err := io.Copy(outfile, file); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if err := a.Insert(); err != nil {
        log.Println(err)
        os.Remove("resource/article/0.md")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    os.Rename("resource/article/0.md", fmt.Sprintf("resource/article/%d.md", a.Id))

    // response home
    handleHome(w, r)
}


