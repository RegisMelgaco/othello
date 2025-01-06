package tcp

import (
	"local/othello/domain/entity"
	"net/http"
)

func (a *App) pass(http.ResponseWriter, *http.Request) {
	a.match.Commit(entity.PassAction{
		Authory: entity.NewAuthor(a.match.Self()),
		Next:    a.match.Opponent(),
	})
}
