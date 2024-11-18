package actconn_test

import (
	"local/othello/domain/entity"
	"local/othello/gateways/tcp/actconn"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Sync(t *testing.T) {
	t.Parallel()

	t.Run("should recieve action", func(t *testing.T) {
		t.Parallel()

		t.Run("given multiple messages", func(t *testing.T) {
			t.Parallel()
			conn1, conn2 := newConnections(t)

			want := entity.MessageAction{
				CreatedAt: date,
			}

			in1, out1 := conn1.Sync()
			in2, out2 := conn2.Sync()

			out1 <- want
			assert.Equal(t, &want, <-in2)

			out2 <- want
			assert.Equal(t, &want, <-in1)

			out1 <- want
			assert.Equal(t, &want, <-in2)
		})

		t.Run("given MessageAction", func(t *testing.T) {
			t.Parallel()
			conn1, conn2 := newConnections(t)

			want := entity.MessageAction{
				Authory:   entity.NewAuthor("Tim Maia"),
				CreatedAt: date,
				Text:      "hummmm",
			}

			_, out := conn1.Sync()
			in, _ := conn2.Sync()

			out <- want
			assert.Equal(t, &want, <-in)
		})

		t.Run("given GiveUpAction", func(t *testing.T) {
			t.Parallel()
			conn1, conn2 := newConnections(t)

			want := entity.GiveUpAction{
				Winner:    timMaia.Author(),
				CreatedAt: date,
			}

			_, out := conn1.Sync()
			in, _ := conn2.Sync()

			out <- want
			assert.Equal(t, &want, <-in)
		})

		t.Run("given PassAction", func(t *testing.T) {
			t.Parallel()
			conn1, conn2 := newConnections(t)

			want := entity.PassAction{
				Next:      timMaia.Author(),
				CreatedAt: date,
			}

			_, out := conn1.Sync()
			in, _ := conn2.Sync()

			out <- want
			assert.Equal(t, &want, <-in)
		})

		t.Run("given RemoveAction", func(t *testing.T) {
			t.Parallel()
			conn1, conn2 := newConnections(t)

			want := entity.RemoveAction{
				Pos: entity.BoardPosition{X: 1, Y: 2},
			}

			_, out := conn1.Sync()
			in, _ := conn2.Sync()

			out <- want
			assert.Equal(t, &want, <-in)
		})

		t.Run("given PlaceAction", func(t *testing.T) {
			t.Parallel()
			conn1, conn2 := newConnections(t)

			want := entity.PlaceAction{
				Pos: entity.BoardPosition{X: 1, Y: 2},
				Val: timMaia.Author(),
			}

			_, out := conn1.Sync()
			in, _ := conn2.Sync()

			out <- want
			assert.Equal(t, &want, <-in)
		})
	})
}

func newConnections(t *testing.T) (*actconn.ActConn, *actconn.ActConn) {
	t.Helper()

	portInt, err := getFreePort()
	require.NoError(t, err)

	port := strconv.Itoa(portInt)

	connChan, errChan := actconn.Listen(port)

	cConn, err := actconn.Dial("", port)
	require.NoError(t, err)

	err = <-errChan
	sConn := <-connChan

	require.NoError(t, err)

	return sConn, cConn
}

// valores para teste
var (
	timMaia       = entity.NewAuthor("Tim Maia")
	robertoCarlos = entity.NewAuthor("Roberto Carlos")
	date          = time.Date(1975, 1, 1, 1, 1, 1, 1, time.Local)
)
