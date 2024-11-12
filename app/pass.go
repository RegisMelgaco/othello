package main

import (
	"local/othello/domain/entity"
	"net/http"
)

func (a *App) pass(http.ResponseWriter, *http.Request) {
	a.match.Commit(entity.PassAction{
		Author: a.match.Self(),
	})
}
