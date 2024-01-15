package network

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTCPClient(t *testing.T) {
	t.Parallel()

	request := "hello server"
	response := "hello client"

	listener, err := net.Listen("tcp", ":10001")
	require.NoError(t, err)

	go func() {
		connection, err := listener.Accept()
		if err != nil {
			return
		}

		buffer := make([]byte, 2048)
		count, err := connection.Read(buffer)
		require.NoError(t, err)
		require.Equal(t, []byte(request), buffer[:count])

		_, err = connection.Write([]byte(fmt.Sprintf("%s\n%s\n", response, EndDelim)))
		require.NoError(t, err)

		defer func() {
			err = connection.Close()
			require.NoError(t, err)
			err = listener.Close()
			require.NoError(t, err)
		}()
	}()

	time.Sleep(100 * time.Millisecond)

	client, err := NewTCPClient("127.0.0.1:10001", time.Minute)
	require.NoError(t, err)

	buffer, err := client.Send([]byte(request))
	require.NoError(t, err)
	require.Equal(t, []byte(response), buffer)

	err = client.Close()
	require.NoError(t, err)
}
