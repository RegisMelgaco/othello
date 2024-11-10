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
	a.match = &entity.Match{
		Players: []entity.PlayerName{
			entity.PlayerName(r.FormValue("self-name")),
			entity.PlayerName(r.FormValue("peer-name")),
		},
		Board: entity.NewBoard(),
	}

	a.match.Chat = []entity.MessageAction{
		{
			Author:    "notas",
			CreatedAt: time.Now(),
			Text:      "Suas peças são marcadas em vermelho e as do oponente em azul.\nPara trocar o valor de uma possição, basta clicar nela até que se obtenha o valor desejado.\nAo clickar em uma posição vazia, é colocada uma peça vermelha, ao clickar em uma vermelha ela é trocada por uma azul, e ao clickar em uma peça azul a peça é removida.",
		},
	}

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

		entity.MessageAction{
			Author:    "rede",
			CreatedAt: time.Now(),
			Text:      msg,
		}.Commit(a.match)

		if err == nil {
			return
		}

		err = a.listen(r.FormValue("self-port"))

		msg = "conectado com sucesso"
		if err != nil {
			msg = fmt.Sprintf("falha ao tentar conectar ao servidor tcp (%s). Aguardando conexão do outro client.", err.Error())
		}

		entity.MessageAction{
			Author:    "rede",
			CreatedAt: time.Now(),
			Text:      msg,
		}.Commit(a.match)
	})
}
