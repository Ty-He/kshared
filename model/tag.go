package model 

import (
    "strings"
    "database/sql"

    "github.com/ty/kshared/conf"
)
type tag struct {
    ts []string  
}

func newTag(total string) *tag {
    // Bad! 
    // t := strings.Split(total, " ")
    t := strings.Fields(total)
    for i := range t {
        t[i] = strings.ToLower(t[i])
    }
    return &tag{
        ts: t,
    }
}

// if err; tx will not rollback
func (t *tag) tryInsertAll(tx *sql.Tx) ([]uint32, error) {
    query := `select id from tag where name = ?;`
    select_stmt, err := tx.Prepare(query)
    if err != nil {
        return nil, err
    }
    defer select_stmt.Close()

    insert := `insert into tag (name) values (?);`
    insert_stmt, err := tx.Prepare(insert)
    if err != nil {
        return nil, err
    }
    defer insert_stmt.Close()

    ids := make([]uint32, len(t.ts))
    for i, name := range t.ts {
        err := select_stmt.QueryRow(name).Scan(&ids[i])
        if err == sql.ErrNoRows {
            // new tag, insert into tag...
            if res, err := insert_stmt.Exec(name); err != nil {
                return nil, err
            } else {
                if insert_id, err := res.LastInsertId(); err != nil {
                    return nil, err
                } else {
                    ids[i] = uint32(insert_id)
                }
            }
        } else if err != nil {
            return nil, err
        } // if err 
    } // for
    
    return ids, nil
}


type Tag string
// Get articles which tags contain t.name
func (t *Tag) GetArticle() ([]*ArticleItem, error) {
    search := strings.ToLower(string(*t))
    query := `select filter_article.id, title, author.name, update_time
    from (
        select article.id, title, article.author_id, update_time 
        from article 
        where article.type != ? and article.id in (
            select article_id
            from article_tag join tag on article_tag.tag_id = tag.id 
            where tag.name = ?
        )
    ) filter_article 
    join author on filter_article.author_id = author.id
    order by update_time desc;`

    rows, err := db.Query(query, conf.Invisible(), search)
    if err != nil {
        return nil, err 
    }
    defer rows.Close()

    return get_items(rows, lastest), nil
}
