package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

const (
	MutexTTL   = 15 * time.Second
	MutexTries = 100
	MutexDelay = 100 * time.Millisecond
)

type Cache struct {
	client  *redis.Client
	redsync *redsync.Redsync
	appName string
}

func New(url, appName string) *Cache {
	redisOptions, err := redis.ParseURL(url)
	if err != nil {
		log.Fatal(err)
	}

	client := redis.NewClient(redisOptions)

	return &Cache{
		client:  client,
		redsync: redsync.New(goredis.NewPool(client)),
		appName: appName,
	}
}

func (c *Cache) Mutex(ctx context.Context, key string) (func(context.Context) (bool, error), error) {
	mutex := c.redsync.NewMutex(c.makeKey("mutex", key),
		redsync.WithExpiry(MutexTTL),
		redsync.WithTries(MutexTries),
		redsync.WithRetryDelay(MutexDelay),
	)

	if err := mutex.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to acquire mutex lock for key %s: %w", key, err)
	}

	return mutex.UnlockContext, nil
}

func (c *Cache) makeKey(prefix, id string) string {
	return fmt.Sprintf("%s:%s:%s", c.appName, prefix, id)
}
