package usecase

import (
	"context"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
	"github.com/Temisaputra/warOnk/internal/repository"
	"github.com/Temisaputra/warOnk/pb/userpb"
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

func (u *UserUsecase) GetAllUsers(ctx context.Context) (res *userpb.GetAllUserResponse, err error) {
	users, err := u.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var userResponses []*userpb.User
	for _, user := range users {
		userResponses = append(userResponses, &userpb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		})
	}

	return &userpb.GetAllUserResponse{Users: userResponses}, nil
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

func (u *UserUsecase) UpdateUser(ctx context.Context, id int32, user *presenter.UserRequest) (err error) {
	updatedUser := entity.Users{
		ID:       id,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Password: user.Password,
	}

	err = u.userRepo.UpdateUser(ctx, updatedUser)
	if err != nil {
		return err
	}

	return
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id int32) (err error) {
	err = u.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return
}
