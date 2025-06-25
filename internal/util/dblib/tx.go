package dblib

import (
	"context"

	"gorm.io/gorm"
)

type TxManager interface {
	DB(ctx context.Context) *gorm.DB
	Tx(ctx context.Context, do func(ctx context.Context) error) error
}

type ContextKeyTx struct{}

type txMan struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) TxManager {
	return &txMan{
		db: db,
	}
}

// DB returns the current transaction if it exists, otherwise it returns the database connection with context.
func (t *txMan) DB(ctx context.Context) *gorm.DB {
	tx, found := t.getCtxTx(ctx)
	if found && tx.Error == nil {
		return tx
	}

	return t.db.WithContext(ctx)
}

// Tx executes a function within a transaction context. If a transaction already exists in the context, it uses that.
func (t *txMan) Tx(ctx context.Context, do func(ctx context.Context) error) error {
	tx, found := t.getCtxTx(ctx)
	if found && tx.Error == nil {
		return do(ctx)
	}

	tx = t.db.Begin()

	if err := do(t.setCtxTx(ctx, tx)); err != nil {
		tx.Rollback()

		return err
	}

	tx.Commit()

	return nil
}

func (t *txMan) setCtxTx(ctx context.Context, tx *gorm.DB) context.Context {
	if tx.Error != nil {
		return ctx
	}

	return context.WithValue(ctx, ContextKeyTx{}, tx)
}

func (t *txMan) getCtxTx(ctx context.Context) (*gorm.DB, bool) {
	tx, found := ctx.Value(ContextKeyTx{}).(*gorm.DB)

	return tx, found
}
