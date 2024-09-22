package model

import (
    "fmt"
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

// if args are invalid, return nil, err; used when don't know id
func NewArticleByItem(atitle, atype, alabel, authorId string) (*Article, error) {
    for _, r := range alabel {
        if r > 255 {
            return nil, errors.New("NewArticle: Invalid label")
        }
    }

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

    var err error
    a.AuthorId, err = parseUint32(authorId)
    if err != nil {
        return nil, err
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

// if existed, update time; else insert and update a.Id
// and if !finished, should renew filesystem
func (a *Article) Insert() error {
    tx, err := db.Begin()
    if err != nil {
        tx.Rollback()
        return err
    }
    query := `insert into article (title, type, label, release_time, update_time, author_id) 
        values (?, ?, ?, ?, ?, ?);`
    result, err := tx.Exec(query, a.Title, a.Type, a.Label, a.UpdateTime, a.UpdateTime, a.AuthorId)
    if err == nil {
        // If in there, insert must ok.
        newId, _ := result.LastInsertId()
        a.Id = uint32(newId)
    } else {
        tx.Rollback()
        return err
    }

    t := newTag(a.Label)
    ids, err := t.tryInsertAll(tx)
    if err != nil {
        tx.Rollback()
        return err
    }

    query = fmt.Sprintf("insert into article_tag (article_id, tag_id) values (%d, ?);", a.Id)
    stmt, err := tx.Prepare(query)
    if err != nil {
        tx.Rollback()
        return err
    }
    defer stmt.Close()

    for _, tag_id := range ids {
        if _, err := stmt.Exec(tag_id); err != nil {
            tx.Rollback()
            return err
        } 
    }

    if err := tx.Commit(); err != nil {
        tx.Rollback()
        return err
    }
    return nil
}

// when know all filed about article 
func NewArticleById(id, authorId string) (*Article, error) {
    a := &Article{
        UpdateTime: time.Now(),
    }
    var err error 
    a.Id, err = parseUint32(id)
    if err != nil {
        return nil, err
    }
    a.AuthorId, err = parseUint32(authorId)
    if err != nil {
        return nil, err
    }
    return a, nil
}

// have a.Id, update new time
func (a *Article) Update() error {
    query := `update article set update_time = ? where id = ? and author_id = ?;`
    result, err := db.Exec(query, a.UpdateTime, a.Id, a.AuthorId)
    if err != nil {
        return err
    } 
    if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
        return errors.New("invalid article id or author_id")
    }
    return nil
}

// Logical delete: type => crash
func (a *Article) Delete() error {
    query := `update article set type = ? where id = ? and author_id = ?;`
    result, err := db.Exec(query, conf.Invisible(), a.Id, a.AuthorId)
    if err != nil {
        return err
    } 
    if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
        return errors.New("invalid article id or author_id")
    }
    return nil
}

// string -> uint32
func parseUint32(s string) (uint32, error) {
    if n, err := strconv.ParseUint(s, 10, 32); err != nil {
        return 0, err
    } else {
        return uint32(n), nil
    }
}
