package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username  string         `json:"username" gorm:"primaryKey"`
	Password  string         `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type LoginInput struct {
	Username string `json:"username" validate:"required,min=1"`
	Password string `json:"password" validate:"required"`
}

func (i *LoginInput) Validate() error {
	return Validator.Struct(i)
}

type UserRepository interface {
	Create(ctx context.Context, u *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
}
