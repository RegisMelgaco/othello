package actconn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"local/othello/domain/entity"
	"log/slog"
	"reflect"
	"time"
)

func (c *ActConn) Send(action entity.Action) error {
	slog.Debug("start conn send", slog.Any("action", action))

	typeName := reflect.TypeOf(action).String()

	enc, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("marshalling action for send: %w", err)
	}

	slog.Debug("conn write", slog.Any("action", action))

	id, err := c.conn.Cmd("%s|%s", typeName, string(enc))
	if err != nil {
		return fmt.Errorf("writing action to connection: %w", err)
	}

	c.conn.StartResponse(id)
	defer c.conn.EndResponse(id)

	if _, err := c.conn.ReadLine(); err != nil {
		return err
	}

	slog.Info("action sent", slog.Any("action", action))

	return nil
}

type ActionListener struct {
	Actions chan entity.Action
	Errs    chan error
}

func (c *ActConn) ListenActions(timeout time.Duration) (ActionListener, error) {
	ln := ActionListener{
		Actions: make(chan entity.Action),
		Errs:    make(chan error),
	}

	go func() {
		for {
			act, err := c.readAction(timeout)
			if err != nil {
				ln.Errs <- err

				continue
			}

			ln.Actions <- act
		}
	}()

	return ln, nil
}

func (c *ActConn) readAction(timeout time.Duration) (entity.Action, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)

	defer cancel()

	var data []byte

	readErr := make(chan error)
	go func() {
		requestID := c.conn.Pipeline.Next()
		c.conn.StartRequest(requestID)
		defer c.conn.EndRequest(requestID)

		var err error
		data, err = c.conn.Reader.ReadLineBytes()
		readErr <- err
	}()

	select {
	case err := <-readErr:
		if err != nil {
			return nil, fmt.Errorf("reading action from connection: %w", err)
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("reading action from connection: %w", ctx.Err())
	}

	slog.Info("msg recieved", slog.String("data", string(data)))

	bs := bytes.Split(data, []byte("|"))
	if len(bs) != 2 {
		return nil, fmt.Errorf("parsing message from socket: %w", fmt.Errorf("unexpected message format"))
	}

	actType, actJSON := bs[0], bs[1]

	ref := reflect.New(actionMap[string(actType)])

	actValue := ref.Interface()

	if err := json.Unmarshal(actJSON, &actValue); err != nil {
		return nil, fmt.Errorf("unmarshaling message from socket: %w", fmt.Errorf("unexpected message format"))
	}

	act, ok := actValue.(entity.Action)
	if !ok {
		return nil, fmt.Errorf("parsing to entity: %w", fmt.Errorf("action could not be parsed to an action"))
	}

	return act, nil
}

var actionMap = map[string]reflect.Type{
	reflect.TypeFor[*entity.MessageAction]().String(): reflect.TypeFor[entity.MessageAction](),
	reflect.TypeFor[*entity.GiveUpAction]().String():  reflect.TypeFor[entity.GiveUpAction](),
	reflect.TypeFor[*entity.PassAction]().String():    reflect.TypeFor[entity.PassAction](),
	reflect.TypeFor[*entity.RemoveAction]().String():  reflect.TypeFor[entity.RemoveAction](),
	reflect.TypeFor[*entity.PlaceAction]().String():   reflect.TypeFor[entity.PlaceAction](),
}
