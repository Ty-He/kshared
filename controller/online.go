package controller 

import (
    "net/http"
    "strconv"
    "log"

    "github.com/ty/kshared/model"
)


func registerOnlineHandle() {
    http.HandleFunc("/login", loginHandle)
    http.HandleFunc("/register", registerHandle)
}

func loginHandle(w http.ResponseWriter, r *http.Request) {
    if (r.Method != http.MethodPost) {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if err := r.ParseForm(); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    id, err := strconv.ParseUint(r.PostFormValue("uid"), 10, 32)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    a := &model.Author{
        Id: uint32(id),
        Pwd: r.PostFormValue("upwd"),
    }

    if err := a.IsValid(); err != nil {
        log.Printf("Author is not valid: Id=%d, Pwd=%s\n Err=%v\n", a.Id, a.Pwd, err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    cookie := http.Cookie{
        Name: "uid",
        Value: r.PostFormValue("uid"),
    }
    http.SetCookie(w, &cookie)

    cookie = http.Cookie{
        Name: "uname",
        Value: a.Name,
    }
    http.SetCookie(w, &cookie)

    // response home page 
    handleHome(w, r)
}

// this register for user, not http router
func registerHandle(w http.ResponseWriter, r *http.Request) {
    if (r.Method != http.MethodPost) {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if err := r.ParseForm(); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    a := &model.Author{
        Name: r.PostFormValue("uname"),
        Pwd: r.PostFormValue("upwd"),
        Email: r.PostFormValue("uemail"),
    }

    if len(a.Name) == 0 {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if err := a.Register(); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    handleHome(w, r)
}


