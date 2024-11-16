package actconn

import (
	"fmt"
	"net"
	"net/textproto"
)

type ActConn struct {
	*textproto.Conn
}

func NewActConn(conn net.Conn) *ActConn {
	return &ActConn{
		Conn: textproto.NewConn(conn),
	}
}

// Dial com um IP vazio, é feito um broadcast
func Dial(ip, port string) (*ActConn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, fmt.Errorf("iniciando conexão tcp: %w", err)
	}

	return NewActConn(conn), nil
}

// Listen usa channels para os retornos para não ser bloqueante na espera por uma conexão com um cliente
func Listen(port string) (connChan chan *ActConn, errChan chan error) {
	connChan = make(chan *ActConn, 1)
	errChan = make(chan error, 1)

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		connChan <- nil
		errChan <- fmt.Errorf("iniciando escuta a porta tcp: %w", err)

		return
	}

	// como ln.Accept é bloqueante, é interessante fazer a chamada desse método em uma routine
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			connChan <- nil
			errChan <- fmt.Errorf("aceitando conexão tcp: %w", err)

			return
		}

		connChan <- NewActConn(conn)
		errChan <- nil
	}()

	return connChan, errChan
}
