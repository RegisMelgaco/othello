package main

import (
	"fmt"
	"local/othello/domain/entity"
	"net/http"
	"slices"
	"time"
)

func (a *App) getChat(w http.ResponseWriter, _ *http.Request) {
	msgs := slices.Clone(a.match.Chat)
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

	entity.MessageAction{
		Author:    a.match.Players[0],
		CreatedAt: time.Now(),
		Text:      msg,
	}.Commit(a.match)

	err := a.templs.chatInput.Execute(w, a.match.Chat)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
