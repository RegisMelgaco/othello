package tcp

import (
	"fmt"
	"local/othello/domain/entity"
	"log/slog"
	"net/http"
	"strconv"
)

func (a *App) updateGame(w http.ResponseWriter, r *http.Request) {
	xStr, yStr := r.URL.Query().Get("pos-x"), r.URL.Query().Get("pos-y")

	x, err := strconv.Atoi(xStr)
	if err != nil {
		slog.Error("pos-x is invalid", slog.String("err", err.Error()), slog.String("pos-x", xStr))
		w.WriteHeader(http.StatusBadRequest)
	}

	y, err := strconv.Atoi(yStr)
	if err != nil {
		slog.Error("pos-y is invalid", slog.String("err", err.Error()), slog.String("pos-y", yStr))
		w.WriteHeader(http.StatusBadRequest)
	}

	act := a.match.HandleClick(entity.BoardPosition{X: x, Y: y})
	slog.Info("board click handled", slog.Any("action", act))

	err = a.templs.board.Execute(w, a.match)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func getGridColors(m entity.Match) [][]string {
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

	slog.Debug("getGridColors", slog.Any("out", out))

	return out
}
