package server

import (
	context "context"

	"github.com/EFG/api"
)

type server struct {
	api.UnimplementedUserServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) GetUsers(ctx context.Context, req *api.GetUsersRequest) (*api.GetUsersResponse, error) {
	// Your implementation here
	return &api.GetUsersResponse{}, nil
}

func (s *server) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	// Your implementation here
	return &api.CreateUserResponse{
		Id:      "123",
		Message: "User created",
	}, nil
}

func (s *server) ModifyUser(ctx context.Context, req *api.ModifyUserRequest) (*api.ModifyUserResponse, error) {
	// Your implementation here
	return &api.ModifyUserResponse{}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	// Your implementation here
	return &api.DeleteUserResponse{}, nil
}
