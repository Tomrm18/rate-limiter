package store

import "errors"

var ErrKeyNotFound = errors.New("key not found")

// Store defines the interface for storing and managing key value pairs
type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
}
