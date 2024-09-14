package model

import (
    "database/sql"
)

type Author struct {
    Id uint32 
    Name string 
    Pwd string 
    Email string // may empty
}

// verify id and pwd, if err != nil, is valid and will fill a.
func (a *Author) IsValid() error {
    query := `select name, email from author where id = ? and pwd = ?;`
    
    row := db.QueryRow(query, a.Id, a.Pwd)
    var s sql.NullString
    if err := row.Scan(&a.Name, &s); err != nil {
        return err
    }
    if s.Valid {
        a.Email = s.String
    }
    return nil
}

// insert into db
func (a *Author) Register() error {
    query := `insert into author (name, pwd, email) values (?, ?, ?)`
    _, err := db.Exec(query, a.Name, a.Pwd, a.Email)
    return err
}
