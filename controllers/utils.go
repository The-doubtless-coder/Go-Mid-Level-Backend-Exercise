package controllers

import (
	"github.com/google/uuid"
)

func ParseUUID(input string) (uuid.UUID, error) {
	return uuid.Parse(input)
}
