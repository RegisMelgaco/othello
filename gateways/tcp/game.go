package tcp

import (
	"fmt"
	"local/othello/domain/entity"
	"net/http"
)

func (a *App) createGame(w http.ResponseWriter, r *http.Request) {
	a.match = entity.NewMatch(
		entity.PlayerName(r.FormValue("self-name")),
		entity.PlayerName(r.FormValue("peer-name")),
	)

	if err := a.connect(r); err != nil {
		a.close(w, err.Error())

		return
	}

	err := a.templs.game.Execute(w, a.match)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
