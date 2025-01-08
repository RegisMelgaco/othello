package grpc_test

import (
	"local/othello/domain/entity"
	"local/othello/gateway/grpc"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Sync(t *testing.T) {
	conn1 := grpc.NewConn(grpc.ConnOpts{
		ListenPort: 3001,
	})

	conn2 := grpc.NewConn(grpc.ConnOpts{
		DialAddress: "127.0.0.1",
		DialPort:    3001,
	})

	want := entity.GiveUpAction{
		Authory: entity.NewAuthor(p2),
		Winner:  p1,
	}

	wait := make(chan bool)
	conn2.OnRecv(func(got entity.Action) {
		assert.Equal(t, want, got)
	})

	conn1.Send(want)

	<-wait
}

const (
	p1 = "tim maia"
	p2 = "roberto carlos"
)
