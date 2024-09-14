package conf 

import (
    "os"
    "log"
    "fmt"
    "encoding/json"
)

// const dsn = "root:245869@tcp(192.168.183.112:3306)/kshared"
type config struct {
    ServerName string `json:"server_name"`
    ServerPort int `json:"server_port"`
    DB struct {
        User string `json:"user"`
        Pwd string `json:"pwd"`
        Ip string `json:"ip"`
        Port int `json:"port"`
    } `json:"db"`
    GroupSize int `json:"group_size"`
    Category []string `json:"category"`
}

var gconfig config 

func init() {
    file, err := os.Open("conf/kshared.json")
    if err != nil {
        log.Fatalln(err)   
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&gconfig); err != nil {
        log.Fatalln(err)   
    }

    log.Println("Load configured file ok!")
}

func ServerName() string {
    return gconfig.ServerName
}

func Port() int {
    return gconfig.ServerPort
}

func GroupSize() int {
    return gconfig.GroupSize 
}

func Dsn() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/kshared", 
        gconfig.DB.User,
        gconfig.DB.Pwd,
        gconfig.DB.Ip,
        gconfig.DB.Port,
    )
}

func Category() []string {
    r := make([]string, len(gconfig.Category))
    copy(r, gconfig.Category)
    return r
}
