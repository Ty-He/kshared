package model

import (
    "time"
    "errors"
    "strconv"

    "github.com/ty/kshared/conf"
)

// for operation, if only display, use ArticleItem
type Article struct {
    Id uint32 
    Title string
    Type string 
    Label string
    UpdateTime time.Time
    AuthorId uint32
}

// if args are invalid, return nil, err
func NewArticle(atitle, atype, alabel, authorId string) (*Article, error) {
    a := &Article{
        Title: atitle,
        Type: atype,
        Label: alabel,
        UpdateTime: time.Now(),
    }
    if len(atitle) == 0 {
        return nil, errors.New("NewArticle: Empty title")
    }
    if !a.isInCategory() {
        return nil, errors.New("NewArticle: Type is invalid")
    }

    if aid, err := strconv.ParseUint(authorId, 10, 32); err != nil {
        return nil, err
    } else {
        a.AuthorId = uint32(aid)
    }
    
    return a, nil 
}

func (a *Article) isInCategory() bool {
    c := conf.Category()
    for i := range c {
        if a.Type == c[i] {
            return true
        }
    }
    return false
}

// if existed, update time; else insert
// and if !finished, should renew filesystem
func (a *Article) Insert() error {
    query := `insert into article (title, type, label, release_time, update_time, author_id) 
        values (?, ?, ?, ?, ?, ?);`
    result, err := db.Exec(query, a.Title, a.Type, a.Label, a.UpdateTime, a.UpdateTime, a.AuthorId)
    if err == nil {
        // If in there, insert must ok.
        newId, _ := result.LastInsertId()
        a.Id = uint32(newId)
    }
    return err
}
