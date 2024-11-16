package tcp

import (
	"fmt"
	"local/othello/domain/entity"
	"net/http"
	"time"
)

func (a *App) pass(http.ResponseWriter, *http.Request) {
	a.match.Commit(entity.PassAction{
		Authory: entity.NewAuthor(a.match.Self()),
		Next:    a.match.Opponent(),
	})

	a.match.Commit(entity.MessageAction{
		Authory:   entity.NewAuthor("jogo"),
		CreatedAt: time.Now(),
		Text:      fmt.Sprintf("<passou a vez para %s>", a.match.Opponent()),
	})
}
