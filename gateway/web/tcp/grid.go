package tcp

import (
	"local/othello/domain/entity"
	"log/slog"
	"net/http"
)

func (a *App) getGrid(w http.ResponseWriter, _ *http.Request) {
	err := a.templs.grid.Execute(w, a.match)
	if err != nil {
		slog.Error("executing grid template", slog.String("err", err.Error()))
	}
}

func getGridColors(m *entity.Match) [][]string {
	out := make([][]string, 8)
	for i, row := range m.Grid() {
		out[i] = make([]string, 8)
		for j, pos := range row {
			switch pos {
			case m.Self():
				out[i][j] = "red"
			case m.Opponent():
				out[i][j] = "blue"
			default:
				out[i][j] = "none"
			}
		}
	}

	return out
}
