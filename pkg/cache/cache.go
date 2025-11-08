package cache

import (
	"context"
	"time"
)

type ICache interface {
	Remember(ctx context.Context, key string, ttl time.Duration, retrieveValueFunc func() (interface{}, error)) ([]byte, error)
	Forget(ctx context.Context, key ...string)
	SetTags(ctx context.Context, key string, tags ...string)
	ForgetTags(ctx context.Context, tags ...string)
	Ping(ctx context.Context) error
}
