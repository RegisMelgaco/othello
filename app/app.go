package main

import (
	"local/othello/domain/entity"
	"net"
)

type App struct {
	match entity.Match
	conn  net.Conn
}
