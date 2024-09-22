package model 

import (
    "database/sql"
)

// source_id send a comment to target_id
type Notify struct {
    Id uint32
    SourceId uint32
    TargetId uint32
    CommentId uint32
    Status string
}

const (
    NotifyStatusUnread = "unread"
    NotifyStatusRead = "read"
    NotifyStatusDeleted = "deleted"
)

// cid is the comment's which sned by this notify target
func newNotify(cid, article_id, sender_id, comment_id uint32) (*Notify, error) {
    n := &Notify{
        SourceId: sender_id,
        CommentId: comment_id,
        Status: NotifyStatusUnread, 
    }
    if (cid == 0) {
        query := `select author_id from article where id = ?;`
        if err := db.QueryRow(query, article_id).Scan(&n.TargetId); err != nil {
            return nil, err
        }
    } else {
        // cid != 0
        query := `select sender_id from comment where id = ?;`
        if err := db.QueryRow(query, cid).Scan(&n.TargetId); err != nil {
            return nil, err
        }
    }
    return n, nil
}

func (n *Notify) insert(tx *sql.Tx) error {
    if n.SourceId == n.TargetId {
        return nil
    }
    query := `insert into notify (source_id, target_id, comment_id, status) values
        (?, ?, ?, ?);`
    if _, err := tx.Exec(query, n.SourceId, n.TargetId, n.CommentId, n.Status); err != nil {
        return err
    }
    return nil
}

// TODO
type NotifyItem struct {
    Id uint32 `json:"id"`
    Sender string `json:"sender"`
    Content string `json:"content"`
    ReleaseTime string `json:"release"`
    ArticleId uint32 `json:"article_id"`
}

func GetUnreadNotifies(uid string) ([]*NotifyItem, error) {
    query := `select cur_notify.id, author.name, content, release_time, article_id
    from (
        select id, source_id, comment_id from notify 
        where target_id = ? and status = ?
    ) cur_notify 
    join author on author.id = source_id 
    join comment on comment.id = comment_id;`
    
    ms := []*NotifyItem{}
    rows, err := db.Query(query, uid, NotifyStatusUnread)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        m := &NotifyItem{}
        rows.Scan(&m.Id, &m.Sender, &m.Content, &m.ReleaseTime, &m.ArticleId)
        ms = append(ms, m)
    }
    return ms, err
}

// TODO 
// func GetReadNotifies(uid string) ([]*Message, error) 
// func GetDeletedNotifies(uid string) ([]*Message, error) 

func MarkedRead(notify_id string) error {
    id, err := parseUint32(notify_id)
    if err != nil {
        return err
    }
    
    query := `update notify set status = ? where id = ?;`
    _, err = db.Exec(query, NotifyStatusRead, id)
    return err
}
