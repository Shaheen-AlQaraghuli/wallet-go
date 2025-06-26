package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"wallet/internal/app/models"
)

const (
	idempotencyKeyPrefix = "idempotency"
	idempotencyTTL       = 24 * time.Hour
)

func (c *Cache) GetIdempotentTransaction(ctx context.Context, idempotencyKey string) (*models.Transaction, error) {
	key := c.makeKey(idempotencyKeyPrefix, idempotencyKey)

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			//nolint:nilnil
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get idempotent transaction from cache: %w", err)
	}

	var transaction models.Transaction
	if err := json.Unmarshal([]byte(val), &transaction); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	return &transaction, nil
}

func (c *Cache) SetIdempotentTransaction(
	ctx context.Context,
	idempotencyKey string,
	transaction models.Transaction,
) error {
	key := c.makeKey(idempotencyKeyPrefix, idempotencyKey)

	data, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	err = c.client.Set(ctx, key, data, idempotencyTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set idempotent transaction in cache: %w", err)
	}

	return nil
}
