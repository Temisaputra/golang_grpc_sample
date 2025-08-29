package repository

import (
	"context"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
	Conn(ctx context.Context) *gorm.DB
}
