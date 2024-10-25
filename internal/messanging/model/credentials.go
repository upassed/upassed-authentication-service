package event

import "github.com/go-playground/validator/v10"

type CredentialsCreateRequest struct {
	Username string `json:"username" validate:"required,min=4,max=30,username"`
	Password string `json:"password" validate:"required,min=8,max=30,password"`
}

func (request *CredentialsCreateRequest) Validate() error {
	validate := validator.New()
	_ = validate.RegisterValidation("username", ValidateUsername())
	_ = validate.RegisterValidation("password", ValidatePassword())

	if err := validate.Struct(*request); err != nil {
		return err
	}

	return nil
}
