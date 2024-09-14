package controller 

import (
    "net/http"
    "os"
    "log"
    "fmt"

    "github.com/ty/kshared/view"
)

func registerMdHandle() {
    http.HandleFunc("/article", handleArticlePage)
    http.HandleFunc("/article_content", handleMd)

    http.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
        content, _ := os.ReadFile("resource/article/6.md")
        w.Write(content)
    })
}
 
// response html page
func handleArticlePage(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    if err := view.ExecuteTmpl(w, &view.TmplArgs{
        Type: "article",
        Value: r.URL.RawQuery,
    }); err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusBadRequest)
        return
    }
}

// reponse a md file from file-system
// this func should not execute template, only response md file for client
func handleMd(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query()["id"]
    if r.Method == http.MethodGet {
        content, err := os.ReadFile(fmt.Sprintf("resource/article/%s.md", id[0]))
        if err != nil {
            log.Println(err)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        w.Write(content)
    } else {
        w.WriteHeader(http.StatusBadRequest)
    }
}

