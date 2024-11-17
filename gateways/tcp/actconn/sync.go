package actconn

import (
	"bytes"
	"encoding/json"
	"local/othello/domain/entity"
	"log/slog"
	"reflect"
)

func (c *ActConn) Sync() (incomming, outgoing chan entity.Action) {
	incomming = make(chan entity.Action, 10)
	outgoing = make(chan entity.Action, 10)

	go func() {
		for {
			c.write(<-outgoing)
		}
	}()

	go func() {
		for {
			if action := c.read(); action != nil {
				incomming <- action
			}
		}
	}()

	return incomming, outgoing
}

func (c *ActConn) write(action entity.Action) {
	typeName := reflect.TypeOf(action).String()

	enc, err := json.Marshal(action)
	if err != nil {
		slog.Error("marshalling action for send", slog.String("err", err.Error()))

		return
	}

	err = c.conn.PrintfLine("%s|%s", typeName, string(enc))
	if err != nil {
		slog.Error("writing action to connection", slog.String("err", err.Error()))

		return
	}

	slog.Info("action sent", slog.Any("action", action))
}

func (c *ActConn) read() entity.Action {
	data, err := c.conn.Reader.ReadLineBytes()
	if err != nil {
		slog.Error("reading from connection", slog.String("err", err.Error()))

		return nil
	}

	slog.Info("msg recieved", slog.String("data", string(data)))

	parts := bytes.Split(data, []byte("|"))
	if len(parts) != 2 && len(parts) != 1 {
		slog.Error("splitting message parts from socket", slog.String("err", "unexpected message format"))

		return nil
	}

	if len(parts) == 1 {
		return nil
	}

	actType, actJSON := parts[0], parts[1]

	ref := reflect.New(actionMap[string(actType)])

	actValue := ref.Interface()

	if err := json.Unmarshal(actJSON, &actValue); err != nil {
		slog.Error("unmarshaling message from socket", slog.String("err", "unexpected message format"))

		return nil
	}

	act, ok := actValue.(entity.Action)
	if !ok {
		slog.Error("parsing to entity", slog.String("err", "action could not be parsed to an action"))

		return nil
	}

	return act
}

var actionMap = map[string]reflect.Type{
	reflect.TypeFor[*entity.MessageAction]().String(): reflect.TypeFor[entity.MessageAction](),
	reflect.TypeFor[*entity.GiveUpAction]().String():  reflect.TypeFor[entity.GiveUpAction](),
	reflect.TypeFor[*entity.PassAction]().String():    reflect.TypeFor[entity.PassAction](),
	reflect.TypeFor[*entity.RemoveAction]().String():  reflect.TypeFor[entity.RemoveAction](),
	reflect.TypeFor[*entity.PlaceAction]().String():   reflect.TypeFor[entity.PlaceAction](),
}
