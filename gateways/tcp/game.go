package tcp

import (
	"fmt"
	"local/othello/domain/entity"
	"net/http"
	"sync"
	"time"
)

var connectOnce sync.Once

func (a *App) getGame(w http.ResponseWriter, r *http.Request) {
	a.match = entity.NewMatch(
		entity.PlayerName(r.FormValue("self-name")),
		entity.PlayerName(r.FormValue("peer-name")),
	)

	err := a.templs.game.Execute(w, a.match)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	go connectOnce.Do(func() {
		err := a.dial(
			r.FormValue("peer-ip"),
			r.FormValue("peer-port"),
		)

		if err != nil {
			a.match.Commit(entity.MessageAction{
				Authory:   entity.NewAuthor("rede"),
				CreatedAt: time.Now(),
				Text:      fmt.Sprintf("falha ao tentar conectar ao servidor tcp (%s). Aguardando pedido de conexão do outro client.", err.Error()),
			})

			err = a.listen(r.FormValue("self-port"))

			if err != nil {
				a.close(w, "não foi possível conectar ao adversário")

				return
			}

			a.match.Commit(entity.PassAction{
				Authory: entity.NewAuthor("server"),
				Next:    a.match.Self(),
			})
		}

		a.match.Commit(entity.MessageAction{
			Authory:   entity.NewAuthor("rede"),
			CreatedAt: time.Now(),
			Text:      "conectado com sucesso",
		})

		a.match.OnCommit(func(act entity.Action) {
			// se a ação não é do próprio jogador, não é necessário envia-la
			if act.Author() != a.match.Self() {
				return
			}

			a.sendAction(act)
		})

		for {
			a.consumeActions()
		}
	})
}
