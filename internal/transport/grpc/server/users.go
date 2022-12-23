package server

import (
	"context"
	pb2 "workshop/internal/transport/grpc/pb"
	"workshop/internal/transport/http/handlers"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

type Users struct {
	pb2.UnimplementedUsersServer
	users handlers.UsersService
}

func NewUsers(us handlers.UsersService) *Users {
	return &Users{users: us}
}

func (u *Users) Create(ctx context.Context, in *pb2.CreateUserRequest) (*pb2.User, error) {
	user, err := u.users.Create(ctx, in.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb2.User{Id: user.ID, Name: user.Name}, nil
}

func (u *Users) GetById(ctx context.Context, in *pb2.GetUserByIdRequest) (*pb2.User, error) {
	user, err := u.users.Get(ctx, in.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb2.User{Id: user.ID, Name: user.Name}, nil
}
