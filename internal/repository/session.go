package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/jatis-redis-training/internal/config"
	"github.com/luckyAkbar/jatis-redis-training/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type sessionRepo struct {
	db     *gorm.DB
	cacher model.Cacher
}

func NewSessionRepo(db *gorm.DB, cacher model.Cacher) model.SessionRepository {
	return &sessionRepo{
		db,
		cacher,
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

	if err := r.setToCache(ctx, session.AccessToken, session, config.DefaultAccessTokenExpiry); err != nil {
		// report error and continue
		logger.Error(err)
	}

	return nil
}

func (r *sessionRepo) FindByAccessToken(ctx context.Context, accessToken string) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"token": utils.Dump(accessToken),
	})

	sessionCache, err := r.findFromCache(ctx, accessToken)
	switch err {
	default:
		logger.Error(err)
	case nil:
		return sessionCache, nil
	}

	session := &model.Session{}
	err = r.db.WithContext(ctx).Model(&model.Session{}).Where("access_token = ?", accessToken).Take(session).Error
	switch err {
	default:
		logger.Error(err)
		return nil, err
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	case nil:
		if err := r.setToCache(ctx, session.AccessToken, session, config.DefaultAccessTokenExpiry); err != nil {
			// report error and continue
			logger.Error(err)
		}
		return session, nil
	}
}

func (r *sessionRepo) findFromCache(ctx context.Context, accessToken string) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"token": utils.Dump(accessToken),
	})

	var session *model.Session
	res, err := r.cacher.Get(ctx, accessToken)
	switch err {
	default:
		logger.Error(err)
		return nil, err
	case redis.Nil:
		return nil, ErrNotFound
	case nil:
		break
	}

	if err := json.Unmarshal([]byte(res), session); err != nil {
		logger.Error(err)
		return nil, err
	}

	return session, nil
}

func (r *sessionRepo) setToCache(ctx context.Context, key string, value *model.Session, expiry time.Duration) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"key":   key,
		"value": value,
	})

	valueStr, err := json.Marshal(value)
	if err != nil {
		logger.Error(err)
		return err
	}

	if err := r.cacher.Set(ctx, key, string(valueStr), expiry); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
