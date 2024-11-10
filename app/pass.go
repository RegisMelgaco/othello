package main

import (
	"local/othello/domain/entity"
	"net/http"
)

func (a *App) pass(http.ResponseWriter, *http.Request) {
	entity.PassAction{
		Author: a.match.Players[0],
	}.Commit(a.match)
}
