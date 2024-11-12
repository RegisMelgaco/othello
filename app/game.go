package main

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

		msg := "conectado com sucesso"
		if err != nil {
			msg = fmt.Sprintf("falha ao tentar conectar ao servidor tcp (%s). Aguardando conexão do outro client.", err.Error())
		}

		a.match.Commit(entity.MessageAction{
			Authory:   entity.Authory{Author: "rede"},
			CreatedAt: time.Now(),
			Text:      msg,
		})

		if err == nil {
			return
		}

		err = a.listen(r.FormValue("self-port"))

		msg = "conectado com sucesso"
		if err != nil {
			msg = fmt.Sprintf("falha ao tentar conectar ao servidor tcp (%s). Aguardando conexão do outro client.", err.Error())
		}

		a.match.Commit(entity.MessageAction{
			Authory:   entity.Authory{Author: "rede"},
			CreatedAt: time.Now(),
			Text:      msg,
		})
	})
}
