package entity

import "time"

type Action interface {
	Commit(*Match)
}

type BoardPosition struct {
	X int
	Y int
}

type PlaceAction struct {
	Author PlayerName
	Pos    BoardPosition
}

func (a PlaceAction) Commit(m *Match) {
	m.Board.grid[a.Pos.X][a.Pos.Y] = a.Author
}

type RemoveAction struct {
	Author PlayerName
	Pos    BoardPosition
}

func (a RemoveAction) Commit(m *Match) {
	m.Board.grid[a.Pos.X][a.Pos.Y] = emptyColor
}

type PassAction struct {
	Author PlayerName
}

func (a PassAction) Commit(m *Match) {
	for _, p := range m.Players {
		if a.Author != p {
			m.TurnOwner = p
		}
	}
}

type GiveUpAction struct {
	Author PlayerName
}

func (a GiveUpAction) Commit(m *Match) {
	for _, p := range m.Players {
		if a.Author != p {
			m.Winner = p
		}
	}
}

type MessageAction struct {
	Author    PlayerName
	CreatedAt time.Time
	Text      string
}

func (a MessageAction) Commit(m *Match) {
	m.Chat = append(m.Chat, a)
}
