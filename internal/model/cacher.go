package model

import (
	"context"
	"time"
)

// Cacher :nodoc:
type Cacher interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, exp time.Duration) error
}
