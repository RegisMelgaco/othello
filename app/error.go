package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func errorMsg(w http.ResponseWriter, err error) {
	templ, _ := template.ParseFS(templates, "templates/error.tmpl.html")
	fmt.Println("err: ", err)

	w.WriteHeader(http.StatusInternalServerError)
	templ.Execute(w, ErrorData{Msg: err.Error()})
}

type ErrorData struct {
	Msg string
}
