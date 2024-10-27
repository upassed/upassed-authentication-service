package domain

import "github.com/google/uuid"

type Credentials struct {
	ID           uuid.UUID
	Username     string
	PasswordHash []byte
}