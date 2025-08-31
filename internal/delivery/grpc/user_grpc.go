package grpc

import (
	"context"

	"github.com/Temisaputra/warOnk/internal/delivery/presenter"
	"github.com/Temisaputra/warOnk/internal/usecase"
	"github.com/Temisaputra/warOnk/pb/userpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer
	uc usecase.UserUsecase
}

func NewUserServiceServer(uc usecase.UserUsecase) *UserServiceServer {
	return &UserServiceServer{uc: uc}
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := s.uc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserResponse{
		User: &userpb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
			// Password: user.Password,
		},
	}, nil
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*emptypb.Empty, error) {
	newUser := &presenter.UserRequest{
		Username: req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: req.Password,
	}

	err := s.uc.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
