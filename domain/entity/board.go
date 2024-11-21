package entity

type Board struct {
	Grid [][]PlayerName
}

const EmptyColor = "none"

func NewBoard() *Board {
	b := Board{make([][]PlayerName, 0, 8)}
	for range 8 {
		row := make([]PlayerName, 0, 8)
		for range 8 {
			row = append(row, EmptyColor)
		}

		b.Grid = append(b.Grid, row)
	}

	return &b
}

func (m *Match) FindWinner() *PlayerName {
	count := make(map[PlayerName]int, 2)
	for _, row := range m.board.Grid {
		for _, v := range row {
			if v == EmptyColor {
				return nil
			}

			if _, ok := count[v]; !ok {
				count[v] = 0
			}

			count[v] += 1
		}
	}

	w := m.self
	if count[m.self] < count[m.opponent] {
		w = m.opponent
	}

	return &w
}

func (m *Match) Grid() [][]PlayerName {
	return m.board.Grid
}
