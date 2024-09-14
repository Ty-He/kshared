package model 

import (
    "log"
    "database/sql"

    _ "github.com/go-sql-driver/mysql"

    "github.com/ty/kshared/conf"
)

var db* sql.DB
// const dsn = "root:245869@tcp(192.168.183.112:3306)/kshared"

// shuould be not motified
var lastest = conf.GroupSize();

func init() {
    var err error
    db, err = sql.Open("mysql", conf.Dsn())
    if err != nil {
        log.Fatalln(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatalln(err)
    }

    log.Println("Conn mysql ok!")
}
