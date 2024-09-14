package view 

import (
    "log"
    "io"
    "text/template"
)

var layout *template.Template

type TmplArgs struct {
    Type string
    Value any
}

func init() {
    layout = template.New("layout")

    var err error
    layout, err = layout.ParseGlob("resource/web/tmpl/*.html")
    if err != nil {
        log.Panicln(err)
    }

    log.Println("Layout template init ok!")
}

// w should be ResponseWriter, for reply web client
func ExecuteTmpl(w io.Writer, args *TmplArgs) error {
    return layout.ExecuteTemplate(w, "layout", args)
}
