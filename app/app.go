package main

import (
	"embed"
	"fmt"
	"html/template"
	"local/othello/domain/entity"
	"net"
)

type App struct {
	match  entity.Match
	conn   net.Conn
	templs struct {
		game      *template.Template
		home      *template.Template
		board     *template.Template
		chatMsgs  *template.Template
		chatInput *template.Template
	}
}

//go:embed templates/*.tmpl.html
var templates embed.FS

func NewApp() (*App, error) {
	app := &App{}

	var err error

	app.templs.home, err = template.ParseFS(templates, "templates/home.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("parsing home: %w", err)
	}

	app.templs.chatMsgs, err = template.ParseFS(templates, "templates/chat_msgs.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("parsing chat: %w", err)
	}

	app.templs.chatInput, err = template.ParseFS(templates, "templates/chat_input.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("parsing chat: %w", err)
	}

	app.templs.game, err = template.New("game.tmpl.html").Funcs(template.FuncMap{
		"getGridColors": getGridColors,
	}).ParseFS(templates, "templates/game.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("parsing game: %w", err)
	}

	app.templs.board, err = template.New("board.tmpl.html").Funcs(template.FuncMap{
		"getGridColors": getGridColors,
	}).ParseFS(templates, "templates/board.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("parsing board: %w", err)
	}

	return app, nil
}
