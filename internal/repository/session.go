package repository

import (
	"context"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/jatis-redis-training/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type sessionRepo struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) model.SessionRepository {
	return &sessionRepo{
		db,
	}
}

func (r *sessionRepo) Create(ctx context.Context, session *model.Session) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":     utils.DumpIncomingContext(ctx),
		"session": utils.Dump(session),
	})

	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (r *sessionRepo) FindByAccessToken(ctx context.Context, accessToken string) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"token": utils.Dump(accessToken),
	})

	session := &model.Session{}

	err := r.db.WithContext(ctx).Model(&model.Session{}).Where("access_token", accessToken).Take(session).Error
	switch err {
	default:
		logger.Error(err)
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		return session, nil
	}
}
