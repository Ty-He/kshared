package model

import (
    "database/sql"
    "github.com/ty/kshared/conf"
)


// for display in html, if want more operation, use Article
type ArticleItem struct {
    Id uint32
    Title string 
    Author string 

    // To be easy, use type of string directly.
    UpdateTime string 
}


func GetLatestArticles() ([]*ArticleItem, error) {
    sql := `select article.id, title, name, update_time
        from article join author on author_id = author.id
        where type != ?
        order by update_time desc
        limit ?;`

    rows, err := db.Query(sql, conf.Invisible(), lastest)
    if err != nil {
        return nil, err 
    }
    defer rows.Close()

    return get_items(rows, lastest), nil
}

func GetTotalArticle() ([]*ArticleItem, error) {
    sql := `select article.id, title, name, update_time
        from article join author on author_id = author.id
        where type != ?
        order by update_time desc;`

    rows, err := db.Query(sql, conf.Invisible())
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    return get_items(rows, 0), err
}

// get a map: category -> items
func FilterArticleByCategory(category []string) (map[string][]*ArticleItem, error) {
    query := `select article.id, title, name, update_time
        from article join author on author_id = author.id
        where type = ?;`
    
    stmt, err := db.Prepare(query);
    if err != nil {
        return nil, err
    }
    defer stmt.Close()
    
    ret := map[string][]*ArticleItem{}
    for _, c := range category {
        rows, err := stmt.Query(c)
        if err != nil {
            return nil, err
        }
        defer rows.Close()
        items := get_items(rows, 0)
        ret[c] = items
    }
    return ret, nil
}

// scan field from db, the n is recommended size
func get_items(rows *sql.Rows, n int) []*ArticleItem {
    items := make([]*ArticleItem, 0, n)
    for rows.Next() {
        a := &ArticleItem{}
        rows.Scan(&a.Id, &a.Title, &a.Author, &a.UpdateTime)
        items = append(items, a)
    }
    return items
}
