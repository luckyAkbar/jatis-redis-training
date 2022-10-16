package repository

import (
	"context"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/jatis-redis-training/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) model.UserRepository {
	return &userRepo{
		db,
	}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":  utils.DumpIncomingContext(ctx),
		"user": utils.Dump(user),
	})

	_, err := r.GetByUsername(ctx, user.Username)
	switch err {
	default:
		logger.Error(err)
		return err
	case nil:
		return ErrAlreadyExists
	case ErrNotFound:
		break
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":      utils.DumpIncomingContext(ctx),
		"username": username,
	})

	user := &model.User{}
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Take(user).Error
	switch err {
	default:
		logger.Error(err)
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return user, nil
	}
}
