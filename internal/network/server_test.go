package network

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"kv_db/pkg/dlog"
)

func TestTCPServer(t *testing.T) {
	t.Parallel()

	request := []byte("hello server\n")
	response := "hello client"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxConnectionsNumber := 2
	idleTimeout := 1 * time.Second

	server, err := NewTCPServer(":20001", maxConnectionsNumber, idleTimeout, dlog.NewNonSlog())
	require.NoError(t, err)

	go func() {
		err := server.HandleQueries(ctx, func(ctx context.Context, buffer []byte) []byte {
			require.Equal(t, request[:len(request)-1], buffer)
			return []byte(response)
		})
		require.NoError(t, err)
	}()

	connection, err := net.Dial("tcp", "localhost:20001")
	require.NoError(t, err)

	_, err = connection.Write(request)
	require.NoError(t, err)

	expected := fmt.Sprintf("%s\n%s\n", response, EndDelim)
	buffer := make([]byte, 1024)
	count, err := connection.Read(buffer)
	require.NoError(t, err)
	require.Equal(t, []byte(expected), buffer[:count])
}
