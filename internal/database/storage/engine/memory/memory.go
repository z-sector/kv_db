package memory

import (
	"context"
	"sync"

	"kv_db/pkg/lock"
)

type HashTable struct {
	mutex sync.RWMutex
	data  map[string]string
}

func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]string),
	}
}

func (s *HashTable) Set(_ context.Context, key string, value string) error {
	lock.WithLock(&s.mutex, func() {
		s.data[key] = value
	})
	return nil
}

func (s *HashTable) Get(_ context.Context, key string) (string, bool, error) {
	var value string
	var ok bool
	lock.WithLock(s.mutex.RLocker(), func() {
		value, ok = s.data[key]
	})
	if !ok {
		return "", false, nil
	}
	return value, true, nil
}

func (s *HashTable) Delete(_ context.Context, key string) error {
	lock.WithLock(&s.mutex, func() {
		delete(s.data, key)
	})
	return nil
}
