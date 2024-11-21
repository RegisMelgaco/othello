package tcp

import (
	"local/othello/domain/entity"
	"net/http"
)

func (a *App) giveUp(http.ResponseWriter, *http.Request) {
	a.match.Commit(entity.GiveUpAction{
		Authory: entity.NewAuthor(a.match.Self()),
		Winner:  a.match.Opponent(),
	})
}
