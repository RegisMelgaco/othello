package tcp

import (
	"fmt"
	"net/http"
)

func (a *App) getHome(w http.ResponseWriter, r *http.Request) {
	data := homeData{
		Msg: r.URL.Query().Get("msg"),
	}

	err := a.templs.home.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

type homeData struct {
	Msg string
}
