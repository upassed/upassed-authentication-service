package business

import "github.com/google/uuid"

type Credentials struct {
	ID          uuid.UUID
	Username    string
	Password    string
	AccountType AccountType
}

type CreateCredentialsResponse struct {
	CreatedCredentialsID uuid.UUID
}
