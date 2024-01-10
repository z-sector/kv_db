package initialization

import (
	"errors"

	"kv_db/config"
	"kv_db/internal/database/storage"
	"kv_db/internal/database/storage/engine/memory"
)

const (
	inMemoryEngine = "in_memory"
)

func CreateEngine(cfg config.EngineConfig) (storage.Engine, error) {
	switch cfg.Type {
	case "":
		return memory.NewHashTable(), nil
	case inMemoryEngine:
		return memory.NewHashTable(), nil
	}

	return nil, errors.New("engine type is incorrect")
}
