package controller 

import (
    "net/http"
    "log"
    "encoding/json"

    "github.com/ty/kshared/model"
)

func registerCommentHandler() {
    http.HandleFunc("/sending_comment", handleNewComment)
    http.HandleFunc("/fetch_comment", handleGetComment)
}

func handleNewComment(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    cookie, err := r.Cookie("uid");
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    } 

    father:= &struct{
        Pid string `json:"pid"`
        Article_id string `json:"article_id"`
        Content string `json:"content"`
    }{}
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(father); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    c, err := model.NewCommentForPost(father.Pid, father.Article_id, cookie.Value, &father.Content)
    if err != nil {
        log.Println(err, father)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if err := c.Insert(); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func handleGetComment(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    c, err := model.NewCommentForGet(r.URL.Query().Get("id"), r.URL.Query().Get("article_id"))
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    cs, err := c.GetNextLevel()
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cs)
}
