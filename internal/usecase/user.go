package usecase

import (
	"context"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
	"github.com/Temisaputra/warOnk/internal/repository"
)

type UserUsecase struct {
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
}

func NewUserUsecase(userRepository repository.UserRepository, transactionRepository repository.TransactionRepository) *UserUsecase {
	return &UserUsecase{
		userRepo:        userRepository,
		transactionRepo: transactionRepository,
	}
}

func (u *UserUsecase) GetById(ctx context.Context, id int32) (res presenter.UserResponse, err error) {
	user, err := u.userRepo.GetUserById(ctx, id)
	if err != nil {
		return presenter.UserResponse{}, err
	}
	return user, nil
}

func (u *UserUsecase) CreateUser(ctx context.Context, user *presenter.UserRequest) (err error) {
	newUser := entity.Users{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Password: user.Password,
	}

	err = u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return err
	}

	return
}
