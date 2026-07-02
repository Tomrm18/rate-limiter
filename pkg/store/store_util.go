package store

import "encoding/json"

func Marshal[T any](v T) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func Unmarshal[T any](s string) (T, error) {
	var v T
	err := json.Unmarshal([]byte(s), &v)
	return v, err
}
