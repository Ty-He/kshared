package controller 

import (
    "net/http"
    "time"
)


// this func is not called, because cilent can delete cookie.
func registerOfflineHandle() {
    http.HandleFunc("/logout", clearCookie)
}

func clearCookie(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // set expire cookie for delete
    for _, cookie := range r.Cookies() {
        http.SetCookie(w, &http.Cookie{
            Name: cookie.Name,
            Value: "",
            Path: cookie.Path,
            Expires: time.Unix(0, 0),
        })
    }

    w.WriteHeader(http.StatusOK)
    // handleHome(w, r)
}
