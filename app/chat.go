package main

import (
	"fmt"
	"local/othello/domain/entity"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

func (a *App) getChat(w http.ResponseWriter, _ *http.Request) {
	msgs := slices.Clone(a.match.Chat())
	slices.Reverse(msgs)
	err := a.templs.chatMsgs.Execute(w, msgs)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (a *App) updateChat(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")

	a.match.Commit(entity.MessageAction{
		Authory:   entity.Authory{Author: a.match.Self()},
		CreatedAt: time.Now(),
		Text:      msg,
	})

	err := a.templs.chatInput.Execute(w, a.match.Chat())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (a *App) errMsg(msg string, err error) {
	slog.Error(msg, slog.String("err", err.Error()))

	a.match.Commit(entity.MessageAction{
		Authory:   entity.Authory{Author: "servidor"},
		CreatedAt: time.Now(),
		Text:      fmt.Errorf("%s: %w", msg, err).Error(),
	})
}
