package server

import (
	"context"
	"testing"
	"time"

	"github.com/EFG/api"
	"github.com/EFG/internal/datasource/database/postgres"
	"github.com/EFG/internal/notifier"
	"github.com/EFG/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_WritesToDataSource(t *testing.T) {
	mockDatasource := &postgres.MockClient{
		UUID: "123e4567-e89b-12d3-a456-426614174000",
	}

	mockNotifier := &notifier.MockNotifier{}

	mockTimeNow := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// Create the gRPC server with the mock
	srv := NewServer(mockDatasource, mockNotifier, mockTimeNow)

	// Mock gRPC request
	req := &api.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
		Country:   "USA",
		Nickname:  "johndoe",
	}

	// Call the gRPC method
	resp, err := srv.CreateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", resp.Id)
	assert.Equal(t, "Successfully created user", resp.Message)

	// Check the mock datasource call history
	assert.Len(t, mockDatasource.GetCallHistory(), 1)
	assert.Equal(t, "John", mockDatasource.GetCallHistory()[0].FirstName.String)
	assert.Equal(t, "Doe", mockDatasource.GetCallHistory()[0].LastName.String)
	assert.NotEqual(t, "password123", mockDatasource.GetCallHistory()[0].Password.String)
	assert.Equal(t, "USA", mockDatasource.GetCallHistory()[0].Country.String)
	assert.Equal(t, "johndoe", mockDatasource.GetCallHistory()[0].Nickname.String)
	assert.Equal(t, "john.doe@example.com", mockDatasource.GetCallHistory()[0].Email.String)

	// check the mock notifier call history
	assert.True(t, mockNotifier.PublishCalled)
	assert.Len(t, mockNotifier.PublishedMessages, 1)
	assert.Contains(t, string(mockNotifier.PublishedMessages[0]), "create")
	assert.Contains(t, string(mockNotifier.PublishedMessages[0]), "123e4567-e89b-12d3-a456-426614174000")
	assert.Contains(t, string(mockNotifier.PublishedMessages[0]), "2025-01-01T00:00:00Z")

	// Reset and check the mock
	mockDatasource.Reset()
	assert.Len(t, mockDatasource.GetCallHistory(), 0)
}

func TestCreateUsers_DataSourceError(t *testing.T) {
	mockDatasource := &postgres.MockClient{
		UUID:              "123e4567-e89b-12d3-a456-426614174000",
		TestRequiresError: true,
	}

	mockNotifier := &notifier.MockNotifier{}

	mockTimeNow := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// Create the gRPC server with the mock
	srv := NewServer(mockDatasource, mockNotifier, mockTimeNow)

	// Mock gRPC request
	req := &api.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
		Country:   "USA",
		Nickname:  "johndoe",
	}

	// Call the gRPC method
	resp, err := srv.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to create user: mock db error for create user")

	// check the mock notifier call history
	assert.False(t, mockNotifier.PublishCalled)
	assert.Len(t, mockNotifier.PublishedMessages, 0)
}

func TestCreateUser_ValidationErrorsForMissingFieldsInRequest(t *testing.T) {
	mockDatasource := &postgres.MockClient{
		UUID: "123e4567-e89b-12d3-a456-426614174000",
	}

	mockNotifier := &notifier.MockNotifier{}

	mockTimeNow := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// Create the gRPC server with the mock
	srv := NewServer(mockDatasource, mockNotifier, mockTimeNow)

	tests := []struct {
		name        string
		modifyReq   func(req *api.CreateUserRequest)
		expectedErr string
	}{
		{
			name: "missing first name",
			modifyReq: func(req *api.CreateUserRequest) {
				req.FirstName = ""
			},
			expectedErr: "FirstName cannot be empty",
		},
		{
			name: "missing last name",
			modifyReq: func(req *api.CreateUserRequest) {
				req.LastName = ""
			},
			expectedErr: "LastName cannot be empty",
		},
		{
			name: "missing email",
			modifyReq: func(req *api.CreateUserRequest) {
				req.Email = ""
			},
			expectedErr: "Email cannot be empty",
		},
		{
			name: "missing password",
			modifyReq: func(req *api.CreateUserRequest) {
				req.Password = ""
			},
			expectedErr: "Password cannot be empty",
		},
		{
			name: "missing country",
			modifyReq: func(req *api.CreateUserRequest) {
				req.Country = ""
			},
			expectedErr: "Country cannot be empty",
		},
		{
			name: "missing nickname",
			modifyReq: func(req *api.CreateUserRequest) {
				req.Nickname = ""
			},
			expectedErr: "NickName cannot be empty",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &api.CreateUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
				Country:   "USA",
				Nickname:  "johndoe",
			}
			tc.modifyReq(req)

			// Call the gRPC method
			resp, err := srv.CreateUser(context.Background(), req)

			assert.Nil(t, resp)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestModifyUser_WritesToDataSource(t *testing.T) {
	mockDatasource := &postgres.MockClient{
		UUID: "123e4567-e89b-12d3-a456-426614174001",
	}

	mockNotifier := &notifier.MockNotifier{}

	mockTimeNow := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// Create the gRPC server with the mock
	srv := NewServer(mockDatasource, mockNotifier, mockTimeNow)

	// Mock gRPC request
	req := &api.ModifyUserRequest{
		Id:        "123e4567-e89b-12d3-a456-426614174001",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
		Country:   "USA",
		Nickname:  "johndoe",
	}

	// Call the gRPC method
	resp, err := srv.ModifyUser(context.Background(), req)
	assert.NoError(t, err)

	assert.Equal(t, "Successfully modified user", resp.Message)

	// Check the mock datasource call history
	assert.Len(t, mockDatasource.GetCallHistory(), 1)
	assert.Equal(t, "John", mockDatasource.GetCallHistory()[0].FirstName.String)
	assert.Equal(t, "Doe", mockDatasource.GetCallHistory()[0].LastName.String)
	assert.NotEqual(t, "password123", mockDatasource.GetCallHistory()[0].Password.String)
	assert.Equal(t, "USA", mockDatasource.GetCallHistory()[0].Country.String)
	assert.Equal(t, "johndoe", mockDatasource.GetCallHistory()[0].Nickname.String)
	assert.Equal(t, "john.doe@example.com", mockDatasource.GetCallHistory()[0].Email.String)

	// check the mock notifier call history
	assert.True(t, mockNotifier.PublishCalled)
	assert.Len(t, mockNotifier.PublishedMessages, 1)
	assert.Contains(t, string(mockNotifier.PublishedMessages[0]), "modify")
	assert.Contains(t, string(mockNotifier.PublishedMessages[0]), "123e4567-e89b-12d3-a456-426614174001")
}

func TestModifyUser_DataSourceError(t *testing.T) {
	mockDatasource := &postgres.MockClient{
		UUID:              "123e4567-e89b-12d3-a456-426614174001",
		TestRequiresError: true,
	}

	mockNotifier := &notifier.MockNotifier{}

	mockTimeNow := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// Create the gRPC server with the mock
	srv := NewServer(mockDatasource, mockNotifier, mockTimeNow)

	// Mock gRPC request
	req := &api.ModifyUserRequest{
		Id:        "123e4567-e89b-12d3-a456-426614174001",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
		Country:   "USA",
		Nickname:  "johndoe",
	}

	// Call the gRPC method
	resp, err := srv.ModifyUser(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to modify user: mock db error for modify user")

	// check the mock notifier call history
	assert.False(t, mockNotifier.PublishCalled)
	assert.Len(t, mockNotifier.PublishedMessages, 0)
}

func TestModifyUser_ValidationErrorsForMissingFieldsInRequest(t *testing.T) {
	mockDatasource := &postgres.MockClient{
		UUID: "123e4567-e89b-12d3-a456-426614174000",
	}
	mockNotifier := &notifier.MockNotifier{}

	mockTimeNow := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// Create the gRPC server with the mock
	srv := NewServer(mockDatasource, mockNotifier, mockTimeNow)

	tests := []struct {
		name        string
		modifyReq   func(req *api.ModifyUserRequest)
		expectedErr *string
	}{
		{
			name: "missing id",
			modifyReq: func(req *api.ModifyUserRequest) {
				req.Id = ""
			},
			expectedErr: utils.Ptr("Id cannot be empty"),
		},
		{
			name: "missing first name",
			modifyReq: func(req *api.ModifyUserRequest) {
				req.FirstName = ""
			},
			expectedErr: nil,
		},
		{
			name: "missing last name",
			modifyReq: func(req *api.ModifyUserRequest) {
				req.LastName = ""
			},
			expectedErr: nil,
		},
		{
			name: "missing email",
			modifyReq: func(req *api.ModifyUserRequest) {
				req.Email = ""
			},
			expectedErr: nil,
		},
		{
			name: "missing password",
			modifyReq: func(req *api.ModifyUserRequest) {
				req.Password = ""
			},
			expectedErr: nil,
		},
		{
			name: "missing country",
			modifyReq: func(req *api.ModifyUserRequest) {
				req.Country = ""
			},
			expectedErr: nil,
		},
		{
			name: "missing nickname",
			modifyReq: func(req *api.ModifyUserRequest) {
				req.Nickname = ""
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &api.ModifyUserRequest{
				Id:        "123e4567-e89b-12d3-a456-426614174001",
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
				Country:   "USA",
				Nickname:  "johndoe",
			}
			tc.modifyReq(req)

			// Call the gRPC method
			resp, err := srv.ModifyUser(context.Background(), req)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			} else {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *tc.expectedErr)
			}
		})
	}
}

func TestDeleteUser_WritesToDataSource(t *testing.T) {
	mockDatasource := &postgres.MockClient{
		UUID: "123e4567-e89b-12d3-a456-426614174001",
	}

	mockNotifier := &notifier.MockNotifier{}

	mockTimeNow := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// Create the gRPC server with the mock
	srv := NewServer(mockDatasource, mockNotifier, mockTimeNow)

	// Mock gRPC request
	req := &api.DeleteUserRequest{
		Id: "123e4567-e89b-12d3-a456-426614174001",
	}

	// Call the gRPC method
	resp, err := srv.DeleteUser(context.Background(), req)
	assert.NoError(t, err)

	assert.Equal(t, "Successfully deleted user", resp.Message)

	// Check the mock datasource call history
	assert.Len(t, mockDatasource.GetCallHistory(), 1)
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174001", mockDatasource.GetCallHistory()[0].ID.String)

	// check the mock notifier call history
	assert.True(t, mockNotifier.PublishCalled)
	assert.Len(t, mockNotifier.PublishedMessages, 1)
	assert.Contains(t, string(mockNotifier.PublishedMessages[0]), "delete")
	assert.Contains(t, string(mockNotifier.PublishedMessages[0]), "123e4567-e89b-12d3-a456-426614174001")
}

func TestDeleteUser_DataSourceError(t *testing.T) {
	mockDatasource := &postgres.MockClient{
		UUID:              "123e4567-e89b-12d3-a456-426614174001",
		TestRequiresError: true,
	}

	mockNotifier := &notifier.MockNotifier{}

	mockTimeNow := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// Create the gRPC server with the mock
	srv := NewServer(mockDatasource, mockNotifier, mockTimeNow)

	// Mock gRPC request
	req := &api.DeleteUserRequest{
		Id: "123e4567-e89b-12d3-a456-426614174001",
	}

	// Call the gRPC method
	resp, err := srv.DeleteUser(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to delete user: mock db error for delete user")

	// check the mock notifier call history
	assert.False(t, mockNotifier.PublishCalled)
	assert.Len(t, mockNotifier.PublishedMessages, 0)
}

func TestDeleteUser_ValidationErrorsForMissingFieldsInRequest(t *testing.T) {
	mockDatasource := &postgres.MockClient{
		UUID: "123e4567-e89b-12d3-a456-426614174000",
	}

	mockNotifier := &notifier.MockNotifier{}

	mockTimeNow := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// Create the gRPC server with the mock
	srv := NewServer(mockDatasource, mockNotifier, mockTimeNow)

	tests := []struct {
		name        string
		modifyReq   func(req *api.DeleteUserRequest)
		expectedErr *string
	}{
		{
			name: "missing id",
			modifyReq: func(req *api.DeleteUserRequest) {
				req.Id = ""
			},
			expectedErr: utils.Ptr("Id cannot be empty"),
		},
		{
			name: "valid id",
			modifyReq: func(req *api.DeleteUserRequest) {
				req.Id = "123e4567-e89b-12d3-a456-426614174001"
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &api.DeleteUserRequest{
				Id: "123e4567-e89b-12d3-a456-426614174001",
			}
			tc.modifyReq(req)

			// Call the gRPC method
			resp, err := srv.DeleteUser(context.Background(), req)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			} else {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *tc.expectedErr)
			}
		})
	}
}

func TestGetUser_ReadsFromDataSource(t *testing.T) {
	mockDatasource := &postgres.MockClient{}

	mockNotifier := &notifier.MockNotifier{}

	mockTimeNow := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// Create the gRPC server with the mock
	srv := NewServer(mockDatasource, mockNotifier, mockTimeNow)

	// Mock gRPC request
	req := &api.GetUsersRequest{
		Page:     1,
		PageSize: 2,
	}

	// Call the gRPC method
	resp, err := srv.GetUsers(context.Background(), req)
	assert.NoError(t, err)

	assert.Len(t, resp.Users, 2)
	assert.Equal(t, "John", resp.Users[0].FirstName)
	assert.Equal(t, "Doe", resp.Users[0].LastName)
	assert.Equal(t, "john.doe@example.com", resp.Users[0].Email)
	assert.Equal(t, "US", resp.Users[0].Country)
	assert.Equal(t, "johndoe", resp.Users[0].Nickname)
	assert.Equal(t, "Jane", resp.Users[1].FirstName)
	assert.Equal(t, "Doe", resp.Users[1].LastName)
	assert.Equal(t, "jane.doe@example.com", resp.Users[1].Email)
	assert.Equal(t, "US", resp.Users[1].Country)
	assert.Equal(t, "janedoe", resp.Users[1].Nickname)
}

func TestGetUsers_DataSourceError(t *testing.T) {
	mockDatasource := &postgres.MockClient{
		TestRequiresError: true,
	}

	mockNotifier := &notifier.MockNotifier{}

	mockTimeNow := func() time.Time {
		return time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	// Create the gRPC server with the mock
	srv := NewServer(mockDatasource, mockNotifier, mockTimeNow)

	// Mock gRPC request
	req := &api.GetUsersRequest{
		Page:     1,
		PageSize: 2,
	}

	// Call the gRPC method
	resp, err := srv.GetUsers(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to get users: mock db error for get users")
}
