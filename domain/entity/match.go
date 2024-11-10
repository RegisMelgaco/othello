package entity

import "log/slog"

type Match struct {
	TurnOwner PlayerName
	Players   []PlayerName
	Actions   []Action
	Board     *Board
	Chat      []MessageAction
	Winner    PlayerName
}

func (m *Match) HandleClick(pos BoardPosition) Action {
	var act Action

	current := m.Board.Grid[pos.X][pos.Y]
	switch current {
	case m.Players[0]:
		act = PlaceAction{
			Authory: Authory{m.Players[0]},
			Val:     m.Players[1],
			Pos:     pos,
		}
	case m.Players[1]:
		act = RemoveAction{
			Authory: Authory{m.Players[0]},
			Pos:     pos,
		}
	default:
		act = PlaceAction{
			Authory: Authory{m.Players[0]},
			Val:     m.Players[0],
			Pos:     pos,
		}
	}

	m.Actions = append(m.Actions, act)

	act.Commit(m)

	slog.Debug("HandleClick", slog.Any("board", m.Board.Grid))

	return act
}
