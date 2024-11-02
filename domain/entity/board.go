package entity

type Board struct {
	grid [][]PlayerName
}

const emptyColor = "none"

func NewBoard() *Board {
	b := Board{make([][]PlayerName, 0, 8)}
	for range b.grid {
		row := make([]PlayerName, 0, 8)
		for range row {
			row = append(row, emptyColor)
		}

		b.grid = append(b.grid, row)
	}

	return &b
}

func (m *Match) FindWinner() *PlayerName {
	count := make(map[PlayerName]int, 2)
	for _, row := range m.Board.grid {
		for _, v := range row {
			if v == emptyColor {
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
