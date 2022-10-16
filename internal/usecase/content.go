package usecase

import (
	"context"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/jatis-redis-training/internal/model"
	"github.com/sirupsen/logrus"
)

type contentUsecase struct {
	contentRepo model.ContentRepository
}

func NewContentUsecase(contentRepo model.ContentRepository) model.ContentUsecase {
	return &contentUsecase{
		contentRepo,
	}
}

func (u *contentUsecase) GetAllMenu(ctx context.Context) ([]model.Menu, error) {
	var menu []model.Menu
	logger := logrus.WithField("ctx", utils.DumpIncomingContext(ctx))

	menu, err := u.contentRepo.GetAllMenu(ctx)
	if err != nil {
		logger.Error(err)
		return menu, ErrInternal
	}

	return menu, nil
}

func (u *contentUsecase) GetAllData(ctx context.Context) ([]model.Data, error) {
	var data []model.Data
	logger := logrus.WithField("ctx", utils.DumpIncomingContext(ctx))

	data, err := u.contentRepo.GetAllData(ctx)
	if err != nil {
		logger.Error(err)
		return data, ErrInternal
	}

	return data, nil
}
