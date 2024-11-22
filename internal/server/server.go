package server

import (
	context "context"
	"log/slog"

	"github.com/EFG/api"
	"github.com/EFG/internal/datasource/dto"
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
	var getUserArgs dto.GetUsersArgs
	getUserArgs.FromAPI(req)

	usersFromDatasource, count, err := service.GetPaginatedUsersList(ctx, s.Datasource, getUserArgs)
	if err != nil {
		return nil, err
	}

	usersForResponse := service.FromDTOToAPI(usersFromDatasource)

	return &api.GetUsersResponse{
		Users:      usersForResponse,
		TotalCount: int32(count),
	}, nil
}

func (s *server) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	if err := validateCreateUserRequest(req); err != nil {
		slog.Error("failed to validate create user request required fields missing", "error", err)
		return nil, err
	}

	user := service.NewUserFromCreateRequest(req)

	id, err := service.FormatNewUserAndPersist(ctx, s.Datasource, user)
	if err != nil {
		return nil, err
	}

	return &api.CreateUserResponse{
		Id:      id,
		Message: "Successfully created user",
	}, nil
}

func (s *server) ModifyUser(ctx context.Context, req *api.ModifyUserRequest) (*api.ModifyUserResponse, error) {
	if err := validateExistingUserRequest(req); err != nil {
		slog.Error("failed to validate modify user request required fields missing", "error", err)
		return nil, err
	}

	user := service.NewUserFromModifyRequest(req)

	err := service.FormatExistingUserAndPersist(ctx, s.Datasource, user)
	if err != nil {
		return nil, err
	}
	return &api.ModifyUserResponse{
		Message: "Successfully modified user",
	}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	if err := validateExistingUserRequest(&api.ModifyUserRequest{Id: req.Id}); err != nil {
		slog.Error("failed to validate delete user request required fields missing", "error", err)
		return nil, err
	}
	err := service.DeleteUserFromDatasource(ctx, s.Datasource, req.Id)
	if err != nil {
		return nil, err
	}

	return &api.DeleteUserResponse{
		Message: "Successfully deleted user",
	}, nil
}
