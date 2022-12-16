package server

import (
	"context"
	pb2 "workshop/internal/transport/grpc/pb"
	"workshop/internal/transport/http/handlers"
	"workshop/internal/users"
)

type Users struct {
	pb2.UnimplementedUsersServer
	user handlers.UsersService
	repo users.Repository
}

func NewUsers(us handlers.UsersService, repo users.Repository) *Users {
	return &Users{user: us, repo: repo}
}

func (u *Users) Create(ctx context.Context, in *pb2.CreateUserRequest) (*pb2.User, error) {
	user, err := u.user.Create(ctx, in.GetName())
	if err != nil {
		return nil, err
	}

	return &pb2.User{Id: user.ID, Name: user.Name}, nil
}

func (u *Users) GetById(ctx context.Context, in *pb2.GetUserByIdRequest) (*pb2.User, error) {
	user, err := u.repo.GetByID(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &pb2.User{Id: user.ID, Name: user.Name}, nil
}
