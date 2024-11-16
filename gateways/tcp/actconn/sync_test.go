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

		t.Run("given MessageAction", func(t *testing.T) {
			t.Parallel()
			conn1, conn2 := newConnections(t)

			want := &entity.MessageAction{
				CreatedAt: date,
			}

			ln, err := conn2.ListenActions(10 * time.Second)
			require.NoError(t, err)

			err = conn1.Send(want)
			require.NoError(t, err)

			select {
			case got := <-ln.Actions:
				assert.Equal(t, want, got)
			case err := <-ln.Errs:
				assert.NoError(t, err)
			}
		})

		t.Run("given GiveUpAction", func(t *testing.T) {
			t.Parallel()
			conn1, conn2 := newConnections(t)

			want := &entity.GiveUpAction{
				Winner:    timMaia.Author(),
				CreatedAt: date,
			}

			err := conn1.Send(want)
			require.NoError(t, err)

			ln, err := conn2.ListenActions(10 * time.Second)
			require.NoError(t, err)

			select {
			case got := <-ln.Actions:
				assert.Equal(t, want, got)
			case err := <-ln.Errs:
				assert.NoError(t, err)
			}
		})

		t.Run("given PassAction", func(t *testing.T) {
			t.Parallel()
			conn1, conn2 := newConnections(t)

			want := &entity.PassAction{
				Next:      timMaia.Author(),
				CreatedAt: date,
			}

			err := conn1.Send(want)
			require.NoError(t, err)

			ln, err := conn2.ListenActions(10 * time.Second)
			require.NoError(t, err)

			select {
			case got := <-ln.Actions:
				assert.Equal(t, want, got)
			case err := <-ln.Errs:
				assert.NoError(t, err)
			}
		})

		t.Run("given RemoveAction", func(t *testing.T) {
			t.Parallel()
			conn1, conn2 := newConnections(t)

			want := &entity.RemoveAction{
				Pos: entity.BoardPosition{X: 1, Y: 2},
			}

			err := conn1.Send(want)
			require.NoError(t, err)

			ln, err := conn2.ListenActions(10 * time.Second)
			require.NoError(t, err)

			select {
			case got := <-ln.Actions:
				assert.Equal(t, want, got)
			case err := <-ln.Errs:
				assert.NoError(t, err)
			}
		})

		t.Run("given PlaceAction", func(t *testing.T) {
			t.Parallel()
			conn1, conn2 := newConnections(t)

			want := &entity.PlaceAction{
				Pos: entity.BoardPosition{X: 1, Y: 2},
				Val: timMaia.Author(),
			}

			err := conn1.Send(want)
			require.NoError(t, err)

			ln, err := conn2.ListenActions(10 * time.Second)
			require.NoError(t, err)

			select {
			case got := <-ln.Actions:
				assert.Equal(t, want, got)
			case err := <-ln.Errs:
				assert.NoError(t, err)
			}
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
