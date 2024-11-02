package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func (a App) connect(selfIP, selfPort, peerIP, peerPort string) error {
	if selfIP < peerIP {
		err := a.listen(selfPort)

		return fmt.Errorf("ouvindo porta tcp: %w", err)
	}

	err := a.dial(peerIP, peerPort)
	if err != nil {
		return fmt.Errorf("connectando ao peer tcp: %w", err)
	}

	return nil
}

func (a App) listen(port string) error {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("iniciando escuta a porta tcp: %w", err)
	}

	a.conn, err = ln.Accept()
	if err != nil {
		return fmt.Errorf("aceitando conexão tcp: %w", err)
	}

	return nil
}

func (a App) dial(ip, port string) error {
	var err error

	a.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return fmt.Errorf("iniciando conexão tcp: %w", err)
	}

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

func requestIPPort(r *http.Request) (ip, port string, err error) {
	ss := strings.Split(r.RemoteAddr, ":")
	if len(ss) != 2 {
		return "", "", fmt.Errorf("unexpected remote addr format")
	}

	return ss[0], ss[1], nil
}
