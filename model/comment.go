package model 

import (
    "time"
    "errors"
)

type Comment struct {
    Id uint32 
    Pid uint32 // pid = 0 => root
    Article_id uint32 
    Sender_id uint32
    Content *string // no copy
    Release time.Time
}

type ResponseComment struct {
    Id string `json:"id"`
    Sender string `json:"sender"`
    Release string `json:"release"`
    Content string `json:"content"`
    Target string `json:"target"`
}


func NewCommentForPost(pid, article_id, sender_id string, content *string) (*Comment, error) {
    c := &Comment{
        Content: content,
        Release: time.Now(),
    }
    var err error

    if c.Pid, err = parseUint32(pid); err != nil {
        return nil, err
    }
    if c.Article_id, err = parseUint32(article_id); err != nil {
        return nil, err
    }
    if c.Sender_id, err = parseUint32(sender_id); err != nil {
        return nil, err
    }

    return c, nil
}

func NewCommentForGet(id, article_id string) (*Comment, error){
    c := &Comment{}
    var err error
    if c.Id, err = parseUint32(id); err != nil {
        return nil, err
    }

    if c.Article_id, err = parseUint32(article_id); err != nil {
        return nil, err
    }
    
    return c, nil
}

// insert self to db
func (c *Comment) Insert() error {
    if !c.valid_ref() {
        return errors.New("invalid pid")
    }
    query := `insert into comment (pid, article_id, sender_id, content, release_time) values 
        (?, ?, ?, ?, ?);`
    _, err := db.Exec(query, c.Pid, c.Article_id, c.Sender_id, *c.Content, c.Release)
    return err
}


func (c *Comment) GetNextLevel() ([]*ResponseComment, error) {
    var target string
    if c.Id == 0 {
        target = "article"
    } else {
        query := `select name from author
            join comment on author.id = comment.sender_id
            where comment.id = ?;`
        row := db.QueryRow(query, c.Id)
        if err := row.Scan(&target); err != nil {
            return nil, err
        }
    }

    query := `select comment.id, name, comment.release_time, content
        from comment join author on sender_id = author.id
        where pid = ? and article_id = ?;`
    rows, err := db.Query(query, c.Id, c.Article_id)
    if err != nil {
        return nil, err
    }

    cs := []*ResponseComment{}
    for rows.Next() {
        rc := &ResponseComment{
            Target: target,
        }
        rows.Scan(&rc.Id, &rc.Sender, &rc.Release, &rc.Content)
        cs = append(cs, rc)
    }
    return cs, nil
}

func (c *Comment) valid_ref() bool {
    if c.Pid == 0 {
        return true
    }
    query := `select 1 from comment where id = ?;`
    row := db.QueryRow(query, c.Pid)
    if err := row.Scan(); err != nil {
        return false
    }
    return true
}
