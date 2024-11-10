package main

import (
	"fmt"
	"net/http"
)

func (a *App) getHome(w http.ResponseWriter, r *http.Request) {
	err := a.templs.home.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
