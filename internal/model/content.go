package model

import (
	"context"
	"time"
)

const (
	DataCacheKey = "DATA"
	MenuCacheKey = "MENU"
)

type Menu struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ParentID int64  `json:"parent_id"`
	Url      string `json:"url"`
}

type Data struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ContentUsecase interface {
	GetAllMenu(ctx context.Context) ([]Menu, error)
	GetAllData(ctx context.Context) ([]Data, error)
}

type ContentRepository interface {
	GetAllMenu(ctx context.Context) ([]Menu, error)
	GetAllData(ctx context.Context) ([]Data, error)
}
