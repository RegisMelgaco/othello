package entity

import (
	"log/slog"
	"time"
)

type Match struct {
	turnOwner PlayerName
	self      PlayerName
	opponent  PlayerName
	actions   []Action
	board     *Board
	chat      []MessageAction
	winner    PlayerName
}

func NewMatch(self, opponent PlayerName) *Match {
	return &Match{
		board:    NewBoard(),
		opponent: opponent,
		self:     self,
		chat: []MessageAction{
			{
				Authory:   Authory{Author: "notas"},
				CreatedAt: time.Now(),
				Text:      "Suas peças são marcadas em vermelho e as do oponente em azul.\nPara trocar o valor de uma possição, basta clicar nela até que se obtenha o valor desejado.\nAo clickar em uma posição vazia, é colocada uma peça vermelha, ao clickar em uma vermelha ela é trocada por uma azul, e ao clickar em uma peça azul a peça é removida.",
			},
		},
	}
}

func (m *Match) Commit(act Action) {
	err := act.commit(m)
	if err != nil {
		return
	}

	m.actions = append(m.actions, act)
}

func (m *Match) HandleClick(pos BoardPosition) Action {
	var act Action

	current := m.board.Grid[pos.X][pos.Y]
	switch current {
	case m.self:
		act = PlaceAction{
			Authory: Authory{m.self},
			Val:     m.opponent,
			Pos:     pos,
		}
	case m.opponent:
		act = RemoveAction{
			Authory: Authory{m.self},
			Pos:     pos,
		}
	default:
		act = PlaceAction{
			Authory: Authory{m.self},
			Val:     m.self,
			Pos:     pos,
		}
	}

	act.commit(m)

	slog.Debug("HandleClick", slog.Any("board", m.board.Grid))

	return act
}

func (m *Match) Self() PlayerName {
	return m.self
}

func (m *Match) Opponent() PlayerName {
	return m.opponent
}

func (m *Match) Chat() []MessageAction {
	return m.chat
}
