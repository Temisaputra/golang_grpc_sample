package grpc

import (
	"context"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/usecase"
	"github.com/Temisaputra/warOnk/pb/authpb"
)

type AuthServiceServer struct {
	authpb.UnimplementedAuthServiceServer
	uc usecase.AuthUsecase
}

func NewAuthServiceServer(uc usecase.AuthUsecase) *AuthServiceServer {
	return &AuthServiceServer{uc: uc}
}

func (s *AuthServiceServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	newUser := &presenter.RegisterRequest{
		Username: req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: req.Password,
	}

	err := s.uc.Register(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &authpb.RegisterResponse{
		User: &authpb.User{
			Name:  req.Name,
			Email: req.Email,
			Role:  req.Role,
			// Password: req.Password,
		},
	}, nil
}
func (s *AuthServiceServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	loginReq := &presenter.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := s.uc.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}

	return &authpb.LoginResponse{
		AccessToken: res.AccessToken,
		User: &authpb.User{
			Id:    res.User.ID,
			Name:  res.User.Name,
			Email: res.User.Email,
			Role:  res.User.Role,
			// Password: res.User.Password,
		},
	}, nil
}
