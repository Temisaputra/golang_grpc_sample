package usecase

import (
	"context"

	"github.com/Temisaputra/warOnk/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/domain/entity"
	"github.com/Temisaputra/warOnk/internal/repository"
	"github.com/Temisaputra/warOnk/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	authRepo        repository.AuthRepository
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
	jwtService      auth.JwtService
}

func NewAuthUsecase(authRepository repository.AuthRepository, userRepository repository.UserRepository, transactionRepository repository.TransactionRepository, jwtService auth.JwtService) *AuthUsecase {
	return &AuthUsecase{
		authRepo:        authRepository,
		userRepo:        userRepository,
		transactionRepo: transactionRepository,
		jwtService:      jwtService,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, params *presenter.RegisterRequest) (err error) {
	return u.transactionRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Mapping ke entity
		user := entity.Users{
			Username: params.Username,
			Email:    params.Email,
			Role:     params.Role,
			Password: string(hash),
		}

		if err := u.userRepo.CreateUser(txCtx, user); err != nil {
			return err
		}
		return nil
	})
}

func (u *AuthUsecase) Login(ctx context.Context, params *presenter.LoginRequest) (res *presenter.LoginResponse, err error) {
	user, err := u.userRepo.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return nil, err
	}

	token, err := u.jwtService.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	users := user.ToPresenter()

	return &presenter.LoginResponse{
		AccessToken: token,
		User:        *users,
	}, nil
}
