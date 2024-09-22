package controller 


import (
    "net/http"
    "log"
    "encoding/json"

    "github.com/ty/kshared/model"
)
func registerNotifyHandle() {
    http.HandleFunc("/get_unread_notify", handleGetUnreadNotifies)
    http.HandleFunc("/marked_notify_read", handleMarkedRead)
}

func handleGetUnreadNotifies(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    cookie, err := r.Cookie("uid")
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    ms, err := model.GetUnreadNotifies(cookie.Value)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // log.Println(ms)
    // w.WriteHeader(http.StatusOK)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ms)
}

func handleMarkedRead(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    id := r.URL.Query().Get("notify_id")
    
    if err := model.MarkedRead(id); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
}
