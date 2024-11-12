package entity

import (
	"fmt"
	"time"
)

type Action interface {
	commit(*Match) error
	author() PlayerName
}

type Authory struct {
	Author PlayerName
}

func (a Authory) author() PlayerName {
	return a.Author
}

type BoardPosition struct {
	X int
	Y int
}

type PlaceAction struct {
	Authory
	Pos BoardPosition
	Val PlayerName
}

func (a PlaceAction) commit(m *Match) error {
	return m.withAuthory(a, func() {
		m.board.Grid[a.Pos.X][a.Pos.Y] = a.Val
	})
}

type RemoveAction struct {
	Authory
	Pos BoardPosition
}

func (a RemoveAction) commit(m *Match) error {
	return m.withAuthory(a, func() {
		m.board.Grid[a.Pos.X][a.Pos.Y] = EmptyColor
	})
}

type PassAction struct {
	Authory
	Author PlayerName
	Next   PlayerName
}

func (a PassAction) commit(m *Match) error {
	m.turnOwner = a.Next

	return MessageAction{
		Authory:   a.Authory,
		CreatedAt: time.Now(),
		Text:      fmt.Sprintf("passou a vez para %s", m.turnOwner),
	}.commit(m)
}

type GiveUpAction struct {
	Authory
	Winner PlayerName
}

func (a GiveUpAction) commit(m *Match) error {
	m.winner = a.Winner

	return MessageAction{
		Authory:   a.Authory,
		CreatedAt: time.Now(),
		Text:      fmt.Sprintf("%s concedeu a vitoria a %s", a.Author, a.Winner),
	}.commit(m)
}

type MessageAction struct {
	Authory
	CreatedAt time.Time
	Text      string
}

func (a MessageAction) commit(m *Match) error {
	m.chat = append(m.chat, a)

	return nil
}

func (m *Match) withAuthory(a Action, f func()) error {
	if a.author() != m.turnOwner {
		return fmt.Errorf("action author (%s) is not the turn owner %s", a.author())
	}

	f()

	return nil
}
