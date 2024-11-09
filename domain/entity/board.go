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
	for _, row := range m.Board.Grid {
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

	var (
		winner      PlayerName
		winnerCount int
	)
	for p, c := range count {
		if c > winnerCount {
			winner = p
		}
	}

	return &winner
}
