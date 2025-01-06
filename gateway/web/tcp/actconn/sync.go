package actconn

import (
	"bytes"
	"context"
	"encoding/json"
	"local/othello/domain/entity"
	"local/othello/gateway/grpc/gen"
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
	var act gen.Action
	switch t := action.(type) {
	case entity.PlaceAction:
		act.Val = &gen.Action_Place{
			Place: &gen.PlaceAction{
				Author: string(t.Author()),
				Position: &gen.Coordinate{
					X: int64(t.Pos.X),
					Y: int64(t.Pos.Y),
				},
				Val: string(t.Val),
			},
		}
	case entity.GiveUpAction:
		act.Val = &gen.Action_GiveUp{
			GiveUp: &gen.GiveUpAction{
				Author: string(t.Author()),
				Winner: string(t.Winner),
			},
		}
	case entity.MessageAction:
		act.Val = &gen.Action_Message{
			Message: &gen.MessageAction{
				Author: string(t.Author()),
				Text:   t.Text,
			},
		}
	case entity.PassAction:
		act.Val = &gen.Action_Pass{
			Pass: &gen.PassAction{
				Author: string(t.Author()),
				Next:   string(t.Next),
			},
		}
	case entity.RemoveAction:
		act.Val = &gen.Action_Remove{
			Remove: &gen.RemoveAction{
				Author: string(t.Author()),
				Position: &gen.Coordinate{
					X: int64(t.Pos.X),
					Y: int64(t.Pos.Y),
				},
			},
		}
	default:
		panic("action not implemented")
	}

	_, err := c.cli.Sync(context.Background(), &act)
	if err != nil {
		slog.Error("grpc client: sync request", slog.String("err", err.Error()))
	}
}

func (c *ActConn) read() entity.Action {
  c.cli.

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
	reflect.TypeFor[entity.MessageAction]().String(): reflect.TypeFor[entity.MessageAction](),
	reflect.TypeFor[entity.GiveUpAction]().String():  reflect.TypeFor[entity.GiveUpAction](),
	reflect.TypeFor[entity.PassAction]().String():    reflect.TypeFor[entity.PassAction](),
	reflect.TypeFor[entity.RemoveAction]().String():  reflect.TypeFor[entity.RemoveAction](),
	reflect.TypeFor[entity.PlaceAction]().String():   reflect.TypeFor[entity.PlaceAction](),
}
