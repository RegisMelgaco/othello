package actconn_test

import (
	"local/othello/gateways/tcp/actconn"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewConn(t *testing.T) {
	portInt, err := getFreePort()
	require.NoError(t, err)

	port := strconv.Itoa(portInt)

	connChan, errChan := actconn.Listen(port)

	cConn, err := actconn.Dial("", port)
	require.NoError(t, err)
	require.NotNil(t, cConn)

	err = <-errChan
	sConn := <-connChan

	require.NoError(t, err)
	require.NotNil(t, sConn)
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()

			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}

	return
}
