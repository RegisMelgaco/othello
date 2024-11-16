package main

import (
	"flag"
	"fmt"
	"local/othello/gateways/tcp"
	"log/slog"
	"os"
)

var serverAddress string

func init() {
	flag.StringVar(&serverAddress, "addr", ":4000", "aplication server address, if host is omitted, the server will listen to all IPs atributed to the host machine")
}

func main() {
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	app, err := tcp.NewApp()
	if err != nil {
		panic(fmt.Errorf("new app: %w", err))
	}

	app.Serve(serverAddress)
}
