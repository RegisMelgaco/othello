package tcp

import (
	"fmt"
	"local/othello/domain/entity"
	"local/othello/gateways/tcp/actconn"
	"log/slog"
	"net/http"
)

func (a *App) connect(r *http.Request) error {
	ip := r.FormValue("ip")
	port := r.FormValue("port")
	isServer := r.FormValue("listen") == "on"

	var (
		conn *actconn.ActConn
		err  error
	)

	if !isServer {
		conn, err = actconn.Dial(ip, port)
		if err != nil {
			return err
		}
	} else {
		connChan, errChan := actconn.Listen(port)

		select {
		case err = <-errChan:
		case conn = <-connChan:
		}

		if err != nil {
			return err
		}
	}

	if isServer {
		a.match.TurnOwner = a.match.Self()
	} else {
		a.match.TurnOwner = a.match.Opponent()
	}

	a.match.Commit(entity.MessageAction{
		Authory: entity.NewAuthor("jogo"),
		Text:    fmt.Sprintf(`jogador "%s" joga primeiro`, a.match.TurnOwner),
	})

	in, out := conn.Sync()

	a.match.OnCommit(func(action entity.Action) {
		slog.Debug("on commit", slog.Any("action", action))
		if action.Author() != a.match.Self() {
			return
		}

		out <- action
	})

	go func() {
		for a.match != nil {
			a.match.Commit(<-in)
		}
	}()

	return nil
}
