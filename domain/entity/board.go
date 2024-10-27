package entity

type Board struct {
	grid [][]string
}

const emptyColor = "none"

func NewBoard() *Board {
	b := Board{make([][]string, 0, 8)}
	for range b.grid {
		row := make([]string, 0, 8)
		for range row {
			row = append(row, emptyColor)
		}

		b.grid = append(b.grid, row)
	}

	return &b
}

func (m *Match) FindWinner() *PlayerName {
	count := make(map[string]int, 2)
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
		winnerColor string
		winnerCount int
	)
	for color, c := range count {
		if c > winnerCount {
			winnerColor = color
		}
	}

	for _, p := range m.Players {
		if p.Color == winnerColor {
			return &p.PlayerName
		}
	}

	return nil
}
