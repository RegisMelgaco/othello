package tcp

import (
	"fmt"
	"local/othello/domain/entity"
	"local/othello/gateway/grpc"
	"log/slog"
	"net/http"
	"strconv"
)

func (a *App) connect(r *http.Request) error {
	ip := r.FormValue("ip")
	port, err := strconv.Atoi(r.FormValue("port"))
	if err != nil {
		panic(err)
	}

	conn := grpc.NewConn(grpc.ConnOpts{
		ListenPort:  port,
		DialAddress: ip,
		DialPort:    port,
	})

	if a.match.Self() < a.match.Opponent() {
		a.match.TurnOwner = a.match.Self()
	} else {
		a.match.TurnOwner = a.match.Opponent()
	}

	a.match.Commit(entity.MessageAction{
		Authory: entity.NewAuthor("jogo"),
		Text:    fmt.Sprintf(`jogador "%s" joga primeiro`, a.match.TurnOwner),
	})

	a.match.OnCommit(func(action entity.Action) {
		slog.Debug("on commit", slog.Any("action", action))
		if action.Author() != a.match.Self() {
			return
		}

		conn.Send(action)
	})

	conn.OnRecv(func(action entity.Action) {
		a.match.Commit(action)
	})

	return nil
}
