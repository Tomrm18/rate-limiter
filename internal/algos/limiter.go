package algos

type Limiter interface {
	// Allow checks if a request is permitted for the given key
	// a key is usually a IP address, API key or user ID
	Allow(key string) (bool, error)
}
