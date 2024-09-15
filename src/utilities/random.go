package utilities

import (
	"github.com/google/uuid"
)

func GenerateRandomID() string {
	id := uuid.New()

	return id.String()
}
