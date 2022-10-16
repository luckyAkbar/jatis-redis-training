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

type contentRepo struct {
	db     *gorm.DB
	cacher model.Cacher
}

func NewContentRepository(db *gorm.DB, cacher model.Cacher) model.ContentRepository {
	return &contentRepo{
		db,
		cacher,
	}
}

func (r *contentRepo) GetAllData(ctx context.Context) ([]model.Data, error) {
	logger := logrus.WithField("ctx", utils.DumpIncomingContext(ctx))

	dataCache, err := r.getDataFromCache(ctx)
	switch err {
	default:
		logger.Error(err)
	case ErrNotFound:
		break
	case nil:
		return dataCache, nil
	}

	var data []model.Data
	if err := r.db.WithContext(ctx).Model(&model.Data{}).Find(&data).Error; err != nil {
		logger.Error(err)
		return data, err
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		logger.Error(err)
	}

	if err := r.setToCache(ctx, model.DataCacheKey, dataByte, config.DefaultDataCacheExpiry); err != nil {
		logger.Error(err)
	}

	return data, nil
}

func (r *contentRepo) GetAllMenu(ctx context.Context) ([]model.Menu, error) {
	logger := logrus.WithField("ctx", utils.DumpIncomingContext(ctx))

	menuCache, err := r.getMenuFromCache(ctx)
	switch err {
	default:
		logger.Error(err)
	case ErrNotFound:
		break
	case nil:
		return menuCache, nil
	}

	var menu []model.Menu
	if err := r.db.WithContext(ctx).Model(&model.Menu{}).Find(&menu).Error; err != nil {
		logger.Error(err)
		return menu, err
	}

	menuByte, err := json.Marshal(menu)
	if err != nil {
		logger.Error(err)
	}

	if err := r.setToCache(ctx, model.MenuCacheKey, menuByte, config.DefaultDataCacheExpiry); err != nil {
		logger.Error(err)
	}

	return menu, nil
}

func (r *contentRepo) getDataFromCache(ctx context.Context) ([]model.Data, error) {
	var data []model.Data
	logger := logrus.WithField("ctx", utils.DumpIncomingContext(ctx))

	res, err := r.cacher.Get(ctx, model.DataCacheKey)
	switch err {
	default:
		logger.Error(err)
		return data, err
	case redis.Nil:
		return data, ErrNotFound
	case nil:
		break
	}

	if err := json.Unmarshal([]byte(res), &data); err != nil {
		logger.Error(err)
		return data, err
	}

	return data, nil
}

func (r *contentRepo) getMenuFromCache(ctx context.Context) ([]model.Menu, error) {
	var menu []model.Menu
	logger := logrus.WithField("ctx", utils.DumpIncomingContext(ctx))

	res, err := r.cacher.Get(ctx, model.MenuCacheKey)
	switch err {
	default:
		logger.Error(err)
		return menu, err
	case redis.Nil:
		return menu, ErrNotFound
	case nil:
		break
	}

	if err := json.Unmarshal([]byte(res), &menu); err != nil {
		logger.Error(err)
		return menu, err
	}

	return menu, nil
}

func (r *contentRepo) setToCache(ctx context.Context, key string, value []byte, expiry time.Duration) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"value": value,
	})

	if err := r.cacher.Set(ctx, key, string(value), expiry); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
