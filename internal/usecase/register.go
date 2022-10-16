package usecase

import (
	"context"
	"time"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/jatis-redis-training/internal/model"
	"github.com/luckyAkbar/jatis-redis-training/internal/repository"
	"github.com/sirupsen/logrus"
)

type registerUsecase struct {
	userRepo model.UserRepository
}

func NewRegisterUsecase(userRepo model.UserRepository) model.RegisterUsecase {
	return &registerUsecase{
		userRepo,
	}
}

func (u *registerUsecase) Register(ctx context.Context, input *model.RegisterInput) (*model.User, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"input": utils.Dump(input),
	})

	if err := input.Validate(); err != nil {
		return nil, ErrBadRequest
	}

	user := &model.User{
		Username:  input.Username,
		Password:  input.Password,
		CreatedAt: time.Now(),
	}

	err := u.userRepo.Create(ctx, user)
	switch err {
	default:
		logger.Error(err)
		return nil, ErrInternal
	case repository.ErrAlreadyExists:
		return nil, ErrAlreadyExists
	case nil:
		return user, nil
	}
}
