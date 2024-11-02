package main

import (
	"fmt"
	"net"
	"net/http"
)

func main() {
	server := http.NewServeMux()

	app := App{}

	server.HandleFunc("GET /", app.home)
	server.HandleFunc("POST /game", app.game)

	ip, err := hostIP()
	if err != nil {
		panic(err)
	}

	fmt.Printf("ip da máquinha: %s\n", ip)

	ln, err := net.Listen("tcp4", fmt.Sprintf("%s:4000", ip))
	if err != nil {
		panic(err)
	}

	fmt.Printf("cliente web em http://%s!\n", ln.Addr())

	(&http.Server{Handler: server}).Serve(ln)
}
