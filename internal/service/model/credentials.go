package business

import "github.com/google/uuid"

type Credentials struct {
	Username string
	Password string
}

type CreateCredentialsResponse struct {
	CredentialsID uuid.UUID
}
