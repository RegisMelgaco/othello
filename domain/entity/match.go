package entity

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Match struct {
	TurnOwner PlayerName
	self      PlayerName
	opponent  PlayerName
	actions   []Action
	board     *Board
	chat      []MessageAction
	onCommit  []func(Action)
	winOnce   sync.Once
}

func NewMatch(self, opponent PlayerName) *Match {
	return &Match{
		board:    NewBoard(),
		opponent: opponent,
		self:     self,
		chat: []MessageAction{
			{
				Authory:   NewAuthor("jogo"),
				CreatedAt: time.Now(),
				Text:      "Suas peças são marcadas em vermelho e as do oponente em azul.\nPara trocar o valor de uma possição, basta clicar nela até que se obtenha o valor desejado.\nAo clickar em uma posição vazia, é colocada uma peça vermelha, ao clickar em uma vermelha ela é trocada por uma azul, e ao clickar em uma peça azul a peça é removida.",
			},
		},
	}
}

func (m *Match) OnCommit(fs ...func(Action)) {
	m.onCommit = append(m.onCommit, fs...)
}

func (m *Match) Commit(act Action) {
	err := act.commit(m)
	if err != nil {
		return
	}

	slog.Info("actions commited", slog.Any("action", act))

	m.actions = append(m.actions, act)

	for _, f := range m.onCommit {
		f(act)
	}

	if w := m.FindWinner(); w != nil {
		m.winOnce.Do(func() {
			m.Commit(MessageAction{
				Authory: NewAuthor("jogo"),
				Text:    fmt.Sprintf(`O vencedor é "%s"`, *w),
			})
		})
	}
}

func (m *Match) HandleClick(pos BoardPosition) Action {
	var act Action

	current := m.board.Grid[pos.X][pos.Y]
	switch current {
	case m.self:
		act = PlaceAction{
			Authory: NewAuthor(m.self),
			Val:     m.opponent,
			Pos:     pos,
		}
	case m.opponent:
		act = RemoveAction{
			Authory: NewAuthor(m.self),
			Pos:     pos,
		}
	default:
		act = PlaceAction{
			Authory: NewAuthor(m.self),
			Val:     m.self,
			Pos:     pos,
		}
	}

	m.Commit(act)

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
