package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (App) home(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFS(templates, "templates/home.tmpl.html")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	ip, _, err := requestIPPort(r)
	if err != nil {
		errorMsg(w, err)

		return
	}

	err = templ.Execute(w, HomeTemplData{SelfIP: ip})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

type HomeTemplData struct {
	SelfIP string
}
