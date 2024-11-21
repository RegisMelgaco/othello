package entity

import (
	"fmt"
	"time"
)

type Action interface {
	Author() PlayerName
	commit(*Match) error
}

type Authory struct {
	Name PlayerName
}

func NewAuthor(p PlayerName) Authory {
	return Authory{p}
}

func (a Authory) Author() PlayerName {
	return a.Name
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
	Next      PlayerName
	CreatedAt time.Time
}

func (a PassAction) commit(m *Match) error {
	m.TurnOwner = a.Next

	m.Commit(MessageAction{
		Authory:   NewAuthor("jogo"),
		CreatedAt: time.Now(),
		Text:      fmt.Sprintf(`"%s" passou a vez para "%s"`, a.Author(), a.Next),
	})

	return nil
}

type GiveUpAction struct {
	Authory
	Winner    PlayerName
	CreatedAt time.Time
}

func (a GiveUpAction) commit(m *Match) error {
	m.Commit(MessageAction{
		Authory:   NewAuthor("jogo"),
		CreatedAt: time.Now(),
		Text:      fmt.Sprintf(`"%s" concedeu a vitoria a "%s"`, a.Author(), a.Winner),
	})

	return nil
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
	if a.Author() != m.TurnOwner {
		return fmt.Errorf("action author (%s) is not the turn owner (%s)", a.Author(), m.TurnOwner)
	}

	f()

	return nil
}
