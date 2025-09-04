package db

import (
	"context"

	"github.com/Temisaputra/warOnk/internal/domain/entity"
	irepository "github.com/Temisaputra/warOnk/internal/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AuthRepository struct {
	*TransactionRepository
}

func NewAuthRepo(db *gorm.DB) irepository.AuthRepository {
	return &AuthRepository{
		TransactionRepository: NewTransactionRepo(db),
	}
}

func (r *AuthRepository) Register(ctx context.Context, user *entity.Users) error {
	if err := r.Conn(ctx).WithContext(ctx).Create(user).Error; err != nil {
		return status.Error(codes.Internal, "failed to register user")
	}
	return nil
}

func (r *AuthRepository) Login(ctx context.Context, email, password string) (*entity.Users, error) {
	var user entity.Users
	if err := r.Conn(ctx).WithContext(ctx).Where("email = ? AND password = ?", email, password).First(&user).Error; err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}
	return &user, nil
}
