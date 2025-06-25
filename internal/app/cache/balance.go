package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	balanceKeyPrefix = "balance:"
	balanceTTL       = 24 * time.Hour
)

func (c *Cache) GetBalance(ctx context.Context, walletID string) (*int, error) {
	key := c.makeKey(balanceKeyPrefix, walletID)

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get balance from cache: %w", err)
	}

	balance, err := strconv.Atoi(val)
	if err != nil {
		return nil, fmt.Errorf("failed to parse balance value: %w", err)
	}

	return &balance, nil
}

func (c *Cache) SetBalance(ctx context.Context, walletID string, balance int) error {
	key := c.makeKey(balanceKeyPrefix, walletID)

	err := c.client.Set(ctx, key, strconv.Itoa(balance), balanceTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set balance in cache: %w", err)
	}

	return nil
}
