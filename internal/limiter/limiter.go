package limiter

import "rate-limiter/internal/common"

type Limiter interface {
	// Allow checks if a request is permitted for the given key
	// a key is usually a IP address, API key or user ID
	Allow(key string) (*common.Result, error)
	// AllowN checks if n requests is permitted for the given key
	// a key is usually a IP address, API key or user ID
	AllowN(key string, n uint) (*common.Result, error)
}
