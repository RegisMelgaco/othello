package main

import (
	"embed"
	"fmt"
	"html/template"
	"local/othello/domain/entity"
	"net/http"
)

//go:embed templates/*
var templates embed.FS

func (a App) game(w http.ResponseWriter, r *http.Request) {
	err := a.connect(
		r.RemoteAddr,
		r.FormValue("self-port"),
		r.FormValue("peer-ip"),
		r.FormValue("peer-port"),
	)

	if err != nil {
		errorMsg(w, err)

		return
	}

	a.match = entity.Match{
		Players: []entity.PlayerName{
			entity.PlayerName(r.FormValue("name")),
		},
		Board: entity.NewBoard(),
	}

	templ, err := template.ParseFS(templates, "templates/game.tmpl.html")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = templ.Execute(w, a.match)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
