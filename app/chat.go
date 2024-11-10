package main

import (
	"fmt"
	"local/othello/domain/entity"
	"net/http"
	"time"
)

func (a *App) getChat(w http.ResponseWriter, _ *http.Request) {
	err := a.templs.chatMsgs.Execute(w, a.match.Chat)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (a *App) updateChat(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	a.match.Chat = append(a.match.Chat, entity.MessageAction{
		Author:    a.match.Players[0],
		CreatedAt: time.Now(),
		Text:      msg,
	})

	err := a.templs.chatInput.Execute(w, a.match.Chat)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
