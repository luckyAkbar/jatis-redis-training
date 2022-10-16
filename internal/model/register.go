package model

import "context"

type RegisterInput struct {
	Username             string `json:"username" validate:"required,min=1"`
	Password             string `json:"password" validate:"required,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

func (r *RegisterInput) Validate() error {
	return Validator.Struct(r)
}

type RegisterUsecase interface {
	Register(ctx context.Context, input *RegisterInput) (*User, error)
}
