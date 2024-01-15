package initialization

import (
	"testing"

	"github.com/stretchr/testify/require"

	"kv_db/config"
)

func TestCreateEngineWithEmptyConfigFields(t *testing.T) {
	t.Parallel()

	engine, err := CreateEngine(config.EngineConfig{})
	require.NoError(t, err)
	require.NotNil(t, engine)
}

func TestCreateEngineWithIncorrectType(t *testing.T) {
	t.Parallel()

	engine, err := CreateEngine(config.EngineConfig{Type: "incorrect"})
	require.Error(t, err)
	require.Nil(t, engine)
}

func TestCreateEngine(t *testing.T) {
	t.Parallel()

	cfg := config.EngineConfig{Type: "in_memory"}

	engine, err := CreateEngine(cfg)
	require.NoError(t, err)
	require.NotNil(t, engine)
}
