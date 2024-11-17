package tcp

import (
	"local/othello/domain/entity"
	"local/othello/gateways/tcp/actconn"
	"log/slog"
	"net/http"
	"time"
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

	a.match.OnCommit(func(action entity.Action) {
		slog.Debug("on commit", slog.Any("action", action))
		if action.Author() != a.match.Self() {
			return
		}

		err := conn.Send(action)
		if err != nil {
			slog.Error("send action", slog.String("err", err.Error()))

			return
		}
	})

	ln, err := conn.ListenActions(5 * time.Minute)
	if err != nil {
		return err
	}

	go func() {
		for a.match != nil {
			select {
			case action := <-ln.Actions:
				a.match.Commit(action)
			case err := <-ln.Errs:
				a.match.Commit(entity.MessageAction{
					Authory:   entity.NewAuthor("servidor"),
					CreatedAt: time.Now(),
					Text:      err.Error(),
				})
			}
		}
	}()

	return nil
}
