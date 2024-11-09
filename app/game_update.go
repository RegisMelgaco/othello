package main

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

	clickedPos := a.match.Board.Grid[x][y]
	switch clickedPos {
	case entity.EmptyColor:
		a.match.Board.Grid[x][y] = a.match.Players[0]
	case a.match.Players[0]:
		a.match.Board.Grid[x][y] = a.match.Players[1]
	case a.match.Players[0]:
		a.match.Board.Grid[x][y] = entity.EmptyColor
	}

	err = a.templs.board.Execute(w, a.match)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func getGridColors(m entity.Match) [][]string {
	out := make([][]string, 8)
	for i, row := range m.Board.Grid {
		out[i] = make([]string, 8)
		for j, pos := range row {
			switch pos {
			case m.Players[0]:
				out[i][j] = "red"
			case m.Players[1]:
				out[i][j] = "blue"
			case m.Players[1]:
				out[i][j] = "none"
			}
		}
	}

	return out
}
