package db

import (
	"context"
	"fmt"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter/response"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// WithTransaction menjalankan function fn dalam satu transaksi
func (r *TransactionRepository) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) (err error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return &response.TxError{Op: "begin", Err: tx.Error}
	}

	// rollback otomatis jika terjadi panic atau context cancel
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			err = &response.TxError{Op: "panic", Err: fmt.Errorf("%v", p)}
		} else if err != nil {
			if rbErr := tx.Rollback().Error; rbErr != nil {
				err = &response.TxError{Op: "rollback", Err: rbErr}
			}
		} else {
			if cErr := tx.Commit().Error; cErr != nil {
				err = &response.TxError{Op: "commit", Err: cErr}
			}
		}
	}()

	// set transaction ke context agar semua repo bisa pakai
	txCtx := context.WithValue(ctx, txContextKey{}, tx)

	// jalankan function user
	err = fn(txCtx)
	return err
}

// Conn mendapatkan DB / transaction dari context
func (r *TransactionRepository) Conn(ctx context.Context) *gorm.DB {
	if tx := ctx.Value(txContextKey{}); tx != nil {
		if txDB, ok := tx.(*gorm.DB); ok {
			return txDB
		}
	}
	return r.db.WithContext(ctx)
}

// key unik untuk context
type txContextKey struct{}
