package main

import (
	"fmt"
	"local/othello/app/static"
	"log/slog"
	"mime"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	mux := chi.NewMux().With(middleware.Logger).With(middleware.Recoverer)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	app, err := NewApp()
	if err != nil {
		panic(fmt.Errorf("new app: %w", err))
	}

	mux.HandleFunc("GET /", app.getHome)
	mux.HandleFunc("GET /game", app.getGame)
	mux.HandleFunc("PUT /game", app.updateGame)
	mux.HandleFunc("GET /chat", app.getChat)
	mux.HandleFunc("PUT /chat", app.updateChat)
	mux.HandleFunc("PUT /pass", app.pass)

	fs := http.FileServerFS(static.FS)
	mux.HandleFunc("GET /static/*", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, ".")
		ext := parts[len(parts)-1]

		w.Header().Add("Content-Type", mime.TypeByExtension("."+ext))

		slog.Debug("mime type set", slog.String("content-type", w.Header().Get("Content-Type")))

		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	})

	ip, err := hostIP()
	if err != nil {
		panic(err)
	}

	fmt.Printf("ip da máquinha: %s\n", ip)

	ln, err := net.Listen("tcp4", ":4000")
	if err != nil {
		panic(err)
	}

	fmt.Printf("cliente web em http://%s!\n", ln.Addr())

	(&http.Server{Handler: mux}).Serve(ln)
}
