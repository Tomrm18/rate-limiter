package algos

import "errors"

var (
	ErrInvalidN             = errors.New("Invalid N input")
	ErrNGreaterThanCapacity = errors.New("Number of requests greater than capacity")
)
