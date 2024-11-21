package server

import (
	context "context"

	"github.com/EFG/api"
	"github.com/EFG/internal/service"
)

type Datasource interface {
	service.Reader
	service.Writer
}

type server struct {
	api.UnimplementedUserServiceServer
	service.Datasource // we are importing the service package anyway seems redundant to define the interface here for package independence
}

func NewServer(d Datasource) *server {
	return &server{
		Datasource: d,
	}
}

func (s *server) GetUsers(ctx context.Context, req *api.GetUsersRequest) (*api.GetUsersResponse, error) {
	return &api.GetUsersResponse{}, nil
}

func (s *server) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	var u service.User
	u.FromAPI(req)

	id, err := service.FormatUserAndPersist(ctx, s.Datasource, u)
	if err != nil {
		return nil, err
	}

	return &api.CreateUserResponse{
		Id:      id,
		Message: "Successfully created user",
	}, nil
}

func (s *server) ModifyUser(ctx context.Context, req *api.ModifyUserRequest) (*api.ModifyUserResponse, error) {
	return &api.ModifyUserResponse{}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	return &api.DeleteUserResponse{}, nil
}
