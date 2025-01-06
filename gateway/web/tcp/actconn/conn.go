package actconn

import (
	"local/othello/gateway/grpc/gen"
)

type ActConn struct {
	cli gen.OthelloClient
}

func NewActConn(cli gen.OthelloClient) *ActConn {
	return &ActConn{cli: cli}
}
