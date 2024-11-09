package main

import (
	"fmt"
	"local/othello/domain/entity"
	"time"
)

func (a *App) errorMsg(err error) {
	a.match.Chat = append(a.match.Chat, entity.MessageAction{
		Author:    "servidor",
		CreatedAt: time.Now(),
		Text:      fmt.Sprintf("erro = %s", err.Error()),
	})

	fmt.Println("err: ", err)
}

type ErrorData struct {
	Err string
}
