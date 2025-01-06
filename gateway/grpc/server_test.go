package grpc_test

import (
	"context"
	"local/othello/domain/entity"
	gateways "local/othello/gateways/grpc"
	"local/othello/gateways/grpc/gen"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/nettest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Test_Server_Place(t *testing.T) {
	t.Run("expect board updated", func(t *testing.T) {
		t.Run("given valid request", func(t *testing.T) {
			match := entity.NewMatch("eu", "cleber")
			match.TurnOwner = entity.PlayerName("cleber")

			s := gateways.NewServer(match)

			ss := grpc.NewServer()
			gen.RegisterOthelloServer(ss, s)

			lis, err := nettest.NewLocalListener("tcp")
			require.NoError(t, err)

			go func() {
				ss.Serve(lis)
			}()

			cc, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
			require.NoError(t, err)

			c := gen.NewOthelloClient(cc)

			_, err = c.Place(context.Background(), &gen.PlaceRequest{
				Player: "cleber",
				Position: &gen.Coordinate{
					X: 1,
					Y: 1,
				},
			})
			require.NoError(t, err)

			assert.Equal(t, entity.PlayerName("cleber"), match.Grid()[1][1])
		})
	})
}
