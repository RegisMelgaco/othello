package grpc

import (
	"context"
	"fmt"
	"io"
	"local/othello/domain/entity"
	"local/othello/gateway/grpc/gen"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type ConnOpts struct {
	DialAddress string
	ListenPort  int
	DialPort    int
}

type ActConn struct {
	in  chan entity.Action
	out chan entity.Action
}

func (c *ActConn) Send(action entity.Action) {
	c.out <- action
}

func (c *ActConn) OnRecv(f func(entity.Action)) {
	go func() {
		for {
			f(<-c.in)
		}
	}()
}

func NewConn(opts ConnOpts) *ActConn {
	conn := &ActConn{
		in:  make(chan entity.Action, 5),
		out: make(chan entity.Action, 5),
	}

	go func() {
		stream, err := dial(opts)
		if err != nil && status.Code(err) != codes.Unavailable {
			panic(err)
		}

		if err != nil {
			stream, err = listen(opts)
			if err != nil {
				panic(err)
			}
		}

		go readActions(stream, conn.in)
		go writeActions(stream, conn.out)
	}()

	return conn
}

type BiStream interface {
	Send(*gen.Action) error
	Recv() (*gen.Action, error)
}

func dial(opts ConnOpts) (BiStream, error) {
	target := fmt.Sprintf("%s:%d", opts.DialAddress, opts.DialPort)
	clientConn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("creating client (target=%s): %w", target, err)
	}

	cli := gen.NewOthelloClient(clientConn)

	stream, err := cli.Sync(context.Background())
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func listen(opts ConnOpts) (BiStream, error) {
	service := NewService()

	addr := fmt.Sprintf(":%d", opts.ListenPort)
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", opts.ListenPort))
		if err != nil {
			panic(err)
		}

		server := grpc.NewServer()
		gen.RegisterOthelloServer(server, service)

		slog.Info("listening", slog.String("addr", addr))
		err = server.Serve(lis)

		panic(err)
	}()

	return <-service.stream, nil
}

func readActions(stream BiStream, in chan entity.Action) {
	for {
		act, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		in <- toEntityAction(act)
	}
}

func writeActions(stream BiStream, out chan entity.Action) {
	for {
		err := stream.Send(toProtoAction(<-out))

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
	}
}

func toEntityAction(action *gen.Action) entity.Action {
	switch t := action.Val.(type) {
	case *gen.Action_Place:
		return entity.PlaceAction{
			Authory: entity.NewAuthor(entity.PlayerName(t.Place.GetAuthor())),
			Pos: entity.BoardPosition{
				X: int(t.Place.GetPosition().GetX()),
				Y: int(t.Place.GetPosition().GetY()),
			},
			Val: entity.PlayerName(t.Place.GetVal()),
		}
	case *gen.Action_GiveUp:
		return entity.GiveUpAction{
			Authory: entity.NewAuthor(entity.PlayerName(t.GiveUp.GetAuthor())),
			Winner:  entity.PlayerName(t.GiveUp.GetWinner()),
		}
	case *gen.Action_Message:
		return entity.MessageAction{
			Authory: entity.NewAuthor(entity.PlayerName(t.Message.GetAuthor())),
			Text:    t.Message.GetText(),
		}
	case *gen.Action_Pass:
		return entity.PassAction{
			Authory: entity.NewAuthor(entity.PlayerName(t.Pass.GetAuthor())),
			Next:    entity.PlayerName(t.Pass.GetNext()),
		}
	case *gen.Action_Remove:
		return entity.RemoveAction{
			Authory: entity.NewAuthor(entity.PlayerName(t.Remove.GetAuthor())),
			Pos: entity.BoardPosition{
				X: int(t.Remove.Position.GetX()),
				Y: int(t.Remove.Position.GetY()),
			},
		}
	}

	panic("action not implemented")
}

func toProtoAction(action entity.Action) *gen.Action {
	act := &gen.Action{}

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

	return act
}
