package tcp

import (
	"fmt"
	"local/othello/gateway/web/tcp/static"
	"log/slog"
	"mime"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) Serve(address string) {
	mux := chi.NewMux().With(middleware.Logger).With(middleware.Recoverer)

	mux.HandleFunc("GET /", a.getHome)
	mux.HandleFunc("GET /game", a.createGame)
	mux.HandleFunc("PUT /game", a.updateGame)
	mux.HandleFunc("GET /grid", a.getGrid)
	mux.HandleFunc("GET /chat", a.getChat)
	mux.HandleFunc("PUT /chat", a.updateChat)
	mux.HandleFunc("PUT /pass", a.pass)
	mux.HandleFunc("PUT /giveup", a.giveUp)

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

	ln, err := net.Listen("tcp4", address)
	if err != nil {
		panic(err)
	}

	fmt.Printf("cliente web em http://%s!\n", ln.Addr())

	(&http.Server{Handler: mux}).Serve(ln)
}

func hostIP() (string, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("acessando tabela de IPs: %w", err)
	}

	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("endereço não encontrado")
}
