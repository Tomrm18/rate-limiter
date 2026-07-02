package mocks

import "github.com/Tomrm18/rate-limiter/pkg/store"

// Mockstore used in testing
type MockStore struct {
	Entries map[string]string
}

func NewMockStore() *MockStore {
	return &MockStore{
		Entries: make(map[string]string),
	}
}

func (m *MockStore) Get(key string) (string, error) {
	if v, ok := m.Entries[key]; !ok {
		return "", store.ErrKeyNotFound
	} else {
		return v, nil
	}
}

func (m *MockStore) Set(key, value string) error {
	m.Entries[key] = value
	return nil
}
