package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"local/othello/domain/entity"
	"net"
	"reflect"
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

	a.match.Commit(entity.PassAction{
		Authory: entity.Authory{Author: "server"},
		Next:    a.match.Self(),
	})

	return nil
}

func (a App) dial(ip, port string) error {
	var err error

	a.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return fmt.Errorf("iniciando conexão tcp: %w", err)
	}

	a.match.Commit(entity.PassAction{
		Author: a.match.Self(),
	})

	return nil
}

func (a App) sendAction(act entity.Action) {
	typeName := []byte(
		reflect.TypeOf(act).Name(),
	)

	enc, err := json.Marshal(act)
	if err != nil {
		a.errMsg("marshalling action to send", err)

		return
	}

	a.conn.Write(bytes.Join([][]byte{typeName, enc}, []byte("|")))
}

func (a App) consumeActions(ctx context.Context) {
	var buff []byte
	_, err := a.conn.Read(buff)
	if err != nil {
		a.errMsg("reading message on socket", err)

		return
	}

	bs := bytes.Split(buff, []byte("|"))
	if len(bs) != 2 {
		a.errMsg("parsing message from socket", fmt.Errorf("unexpected message format"))

		return
	}

	actType, actJSON := bs[0], bs[1]

	actValue := reflect.New(actionMap[string(actType)]).Interface()

	if err := json.Unmarshal(actJSON, &actValue); err != nil {
		a.errMsg("unmarshaling message from socket", fmt.Errorf("unexpected message format"))

		return
	}

	act, ok := actValue.(entity.Action)
	if !ok {
		a.errMsg("parsing to entity", fmt.Errorf("action could not be parsed to an action"))

		return
	}

	a.match.Commit(act)
}

var actionMap = map[string]reflect.Type{
	reflect.TypeFor[entity.MessageAction]().Name(): reflect.TypeFor[entity.MessageAction](),
	reflect.TypeFor[entity.GiveUpAction]().Name():  reflect.TypeFor[entity.GiveUpAction](),
	reflect.TypeFor[entity.PassAction]().Name():    reflect.TypeFor[entity.PassAction](),
	reflect.TypeFor[entity.RemoveAction]().Name():  reflect.TypeFor[entity.RemoveAction](),
	reflect.TypeFor[entity.PlaceAction]().Name():   reflect.TypeFor[entity.PlaceAction](),
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
