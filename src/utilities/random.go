package utilities

import (
	"math/rand"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateRandomID() string {
	id := make([]byte, 12)

	for i := 0; i < len(id); i++ {
		id[i] = charset[rand.Intn(len(charset))]
	}

	return string(id)
}
