package grpc

import (
	"context"
	"local/othello/domain/entity"
	"local/othello/gateways/grpc/gen"
)

type Server struct {
	gen.UnimplementedOthelloServer
	match *entity.Match
}

func NewServer(m *entity.Match) *Server {
	return &Server{match: m}
}

func (s *Server) Place(_ context.Context, req *gen.PlaceRequest) (*gen.Empty, error) {
	s.match.Commit(entity.PlaceAction{
		Authory: entity.NewAuthor(entity.PlayerName(req.GetAuthor())),
		Pos: entity.BoardPosition{
			X: int(req.GetPosition().GetX()),
			Y: int(req.GetPosition().GetY()),
		},
		Val: entity.PlayerName(req.GetVal()),
	})

	return nil, nil
}

func (s *Server) Remove(_ context.Context, req *gen.RemoveRequest) (*gen.Empty, error) {
	s.match.Commit(entity.RemoveAction{
		Authory: entity.NewAuthor(entity.PlayerName(req.GetAuthor())),
		Pos: entity.BoardPosition{
			X: int(req.GetPosition().X),
			Y: int(req.GetPosition().Y),
		},
	})

	return nil, nil
}

func (s *Server) Pass(_ context.Context, req *gen.PassRequest) (*gen.Empty, error) {
	s.match.Commit(entity.PassAction{
		Authory: entity.NewAuthor(entity.PlayerName(req.GetAuthor())),
		Next:    entity.PlayerName(req.Next),
	})

	return nil, nil
}

func (s *Server) GiveUp(_ context.Context, req *gen.GiveUpRequest) (*gen.Empty, error) {
	s.match.Commit(entity.GiveUpAction{
		Authory: entity.NewAuthor(entity.PlayerName(req.GetAuthor())),
		Winner:  entity.PlayerName(req.GetWinner()),
	})

	return nil, nil
}

func (s *Server) Message(_ context.Context, req *gen.MessageRequest) (*gen.Empty, error) {
	s.match.Commit(entity.MessageAction{
		Authory: entity.NewAuthor(entity.PlayerName(req.GetAuthor())),
		Text:    req.GetText(),
	})

	return nil, nil
}
