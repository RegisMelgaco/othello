package entity

import (
	"fmt"
	"time"
)

type Action interface {
	Commit(*Match)
	HasTurn(*Match) bool
}

type Authory struct {
	Author PlayerName
}

func (a Authory) HasTurn(m *Match) bool {
	return a.Author == m.TurnOwner
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

func (a PlaceAction) Commit(m *Match) {
	m.IfHasTurn(a, func() {
		m.Board.Grid[a.Pos.X][a.Pos.Y] = a.Val
	})
}

type RemoveAction struct {
	Authory
	Pos BoardPosition
}

func (a RemoveAction) Commit(m *Match) {
	m.IfHasTurn(a, func() {
		m.Board.Grid[a.Pos.X][a.Pos.Y] = EmptyColor
	})
}

type PassAction struct {
	Authory
	Author PlayerName
	Next   PlayerName
}

func (a PassAction) Commit(m *Match) {
	m.TurnOwner = a.Next

	MessageAction{
		Authory:   a.Authory,
		CreatedAt: time.Now(),
		Text:      fmt.Sprintf("passou a vez para %s", m.TurnOwner),
	}.Commit(m)
}

type GiveUpAction struct {
	Authory
	Winner PlayerName
}

func (a GiveUpAction) Commit(m *Match) {
	m.Winner = a.Winner
}

type MessageAction struct {
	Authory
	CreatedAt time.Time
	Text      string
}

func (a MessageAction) Commit(m *Match) {
	m.Chat = append(m.Chat, a)
}

func (m *Match) IfHasTurn(a Action, f func()) {
	if a.HasTurn(m) {
		f()
	}
}
