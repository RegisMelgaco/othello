package main

import (
	"fmt"
	"local/othello/domain/entity"
	"net"
)

func (a App) listen(port string) error {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("iniciando escuta a porta tcp: %w", err)
	}

	a.conn, err = ln.Accept()
	if err != nil {
		return fmt.Errorf("aceitando conexão tcp: %w", err)
	}

	a.match.TurnOwner = a.match.Players[0]

	return nil
}

func (a App) dial(ip, port string) error {
	var err error

	a.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return fmt.Errorf("iniciando conexão tcp: %w", err)
	}

	entity.PassAction{
		Author: a.match.Players[0],
	}.Commit(a.match)

	return nil
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
