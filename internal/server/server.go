package server

import (
	context "context"
	"log/slog"
	"time"

	"github.com/EFG/api"
	"github.com/EFG/internal/datasource/dto"
	"github.com/EFG/internal/service"
)

type server struct {
	api.UnimplementedUserServiceServer
	// we are importing the service package anyway seems redundant to define the interface here for package independence
	service.Datasource
	service.Notifier
	timeNow func() time.Time
}

func NewServer(d service.Datasource, n service.Notifier, tn func() time.Time) *server {
	return &server{
		Datasource: d,
		Notifier:   n,
		timeNow:    tn,
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

	userChangeNotification := service.CreateUserChangeNotification("create", id, s.timeNow())

	err = service.NotifyOfUserChange(ctx, s.Notifier, userChangeNotification)
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

	userChangeNotification := service.CreateUserChangeNotification("modify", user.ID, s.timeNow())

	err = service.NotifyOfUserChange(ctx, s.Notifier, userChangeNotification)
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

	userChangeNotification := service.CreateUserChangeNotification("delete", req.Id, s.timeNow())

	err = service.NotifyOfUserChange(ctx, s.Notifier, userChangeNotification)
	if err != nil {
		return nil, err
	}

	return &api.DeleteUserResponse{
		Message: "Successfully deleted user",
	}, nil
}
