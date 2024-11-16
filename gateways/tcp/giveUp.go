package tcp

import (
	"fmt"
	"local/othello/domain/entity"
	"net/http"
	"time"
)

func (a *App) giveUp(http.ResponseWriter, *http.Request) {
	a.match.Commit(entity.GiveUpAction{
		Authory: entity.NewAuthor(a.match.Self()),
	})

	a.match.Commit(entity.MessageAction{
		Authory:   entity.NewAuthor("jogo"),
		CreatedAt: time.Now(),
		Text:      fmt.Sprintf("<%s concedeu a vitória a %s>", a.match.Self(), a.match.Opponent()),
	})
}
