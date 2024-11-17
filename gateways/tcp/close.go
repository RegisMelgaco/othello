package tcp

import (
	"fmt"
	"net/http"
)

func (a *App) close(w http.ResponseWriter, msg string) {
	if msg != "" {
		msg = fmt.Sprintf("jogo encerrado: %s", msg)
	}

	w.Header().Add("Location", fmt.Sprintf(`/?msg="%s"`, msg))
	w.WriteHeader(http.StatusSeeOther)
}
