package initialization

import (
	"testing"

	"github.com/stretchr/testify/require"

	"kv_db/config"
	"kv_db/pkg/dlog"
)

func TestCreateNetworkWithEmptyConfigFields(t *testing.T) {
	server, err := CreateNetwork(config.NetworkConfig{}, dlog.NewNonSlog())
	require.NoError(t, err)
	require.NotNil(t, server)
}

func TestCreateNetwork(t *testing.T) {
	t.Parallel()

	cfg := config.NetworkConfig{
		Address:        "localhost:9898",
		MaxConnections: 50,
	}

	server, err := CreateNetwork(cfg, dlog.NewNonSlog())
	require.NoError(t, err)
	require.NotNil(t, server)
}
