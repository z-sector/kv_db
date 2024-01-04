package memory

import "context"

type HashTable struct {
	data map[string]string
}

func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]string),
	}
}

func (s *HashTable) Set(_ context.Context, key string, value string) error {
	s.data[key] = value
	return nil
}

func (s *HashTable) Get(_ context.Context, key string) (string, bool, error) {
	value, ok := s.data[key]
	if !ok {
		return "", false, nil
	}
	return value, true, nil
}

func (s *HashTable) Delete(_ context.Context, key string) error {
	delete(s.data, key)
	return nil
}
