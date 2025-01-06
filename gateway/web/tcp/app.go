package tcp

import (
	"embed"
	"fmt"
	"html/template"
	"local/othello/domain/entity"
)

type App struct {
	match  *entity.Match
	templs struct {
		home      *template.Template
		game      *template.Template
		grid      *template.Template
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
		return nil, fmt.Errorf("parsing home template: %w", err)
	}

	app.templs.chatMsgs, err = template.ParseFS(templates, "templates/chat_msgs.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("parsing chat msgs template: %w", err)
	}

	app.templs.chatInput, err = template.ParseFS(templates, "templates/chat_input.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("parsing chat input template: %w", err)
	}

	app.templs.game, err = template.ParseFS(templates, "templates/game.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("parsing game template: %w", err)
	}

	app.templs.board, err = template.New("board.tmpl.html").Funcs(template.FuncMap{
		"getGridColors": getGridColors,
	}).ParseFS(templates, "templates/board.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("parsing board template: %w", err)
	}

	app.templs.grid, err = template.New("grid.tmpl.html").Funcs(template.FuncMap{
		"getGridColors": getGridColors,
	}).ParseFS(templates, "templates/grid.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("parsing grid template: %w", err)
	}

	return app, nil
}
