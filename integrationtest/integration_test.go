package integrationtest

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/EFG/api"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserIntegration_HappyPath(t *testing.T) {
	client, conn, err := setupGRPCClient("localhost:9000")
	assert.NoError(t, err, "failed to set up gRPC client")
	defer conn.Close()

	d, err := setupPostgresDatasource()
	assert.NoError(t, err, "failed to connect to datasource")
	defer d.Disconnect()

	d.ResetUserStore()

	req := &api.CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "password123",
		Country:   "US",
		Nickname:  "janedoe",
	}

	resp, err := client.CreateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Id, "response ID should not be empty")

	log.Print("Created user with ID: ", resp.Id)

	user, err := d.GetUserById(resp.Id)
	assert.NoError(t, err)

	assert.NoError(t, err)

	assert.Equal(t, req.FirstName, user.FirstName.String)
	assert.Equal(t, req.LastName, user.LastName.String)
	assert.Equal(t, req.Email, user.Email.String)
	assert.Equal(t, req.Country, user.Country.String)
	assert.Equal(t, req.Nickname, user.Nickname.String)
	assert.NotEmpty(t, user.Password.String)
	assert.NotEqual(t, req.Password, user.Password.String)
}

func TestCreateUserIntegration_ErrorCreatingUserForExistingEmail(t *testing.T) {
	client, conn, err := setupGRPCClient("localhost:9000")
	assert.NoError(t, err, "failed to set up gRPC client")
	defer conn.Close()

	d, err := setupPostgresDatasource()
	assert.NoError(t, err, "failed to connect to datasource")
	defer d.Disconnect()

	d.ResetUserStore()

	req := &api.CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "password123",
		Country:   "US",
		Nickname:  "janedoe",
	}

	resp, err := client.CreateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Id, "response ID should not be empty")

	req2 := &api.CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "password123",
		Country:   "US",
		Nickname:  "janedoe",
	}

	resp, err = client.CreateUser(context.Background(), req2)
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Contains(t, err.Error(), "email already exists: jane.doe@example.com")
}

func TestModifyUserIntegration_HappyPath(t *testing.T) {
	client, conn, err := setupGRPCClient("localhost:9000")
	assert.NoError(t, err, "failed to set up gRPC client")
	defer conn.Close()

	d, err := setupPostgresDatasource()
	assert.NoError(t, err, "failed to connect to datasource")
	defer d.Disconnect()

	d.ResetUserStore()

	req := &api.CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "password123",
		Country:   "US",
		Nickname:  "janedoe",
	}

	resp, err := client.CreateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Id, "response ID should not be empty")

	log.Print("Created user with ID: ", resp.Id)

	userBeforeChange, err := d.GetUserById(resp.Id)
	assert.NoError(t, err)

	req2 := &api.ModifyUserRequest{
		Id:        resp.Id,
		FirstName: "John",
		LastName:  "NoDoe",
		Email:     "john.nodoe@example.com",
		Country:   "UK",
		Nickname:  "johnnodoe",
	}

	resp2, err := client.ModifyUser(context.Background(), req2)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully modified user", resp2.Message)

	userAfterChange, err := d.GetUserById(resp.Id)
	assert.NoError(t, err)

	assert.Equal(t, req2.FirstName, userAfterChange.FirstName.String)
	assert.Equal(t, req2.LastName, userAfterChange.LastName.String)
	assert.Equal(t, req2.Email, userAfterChange.Email.String)
	assert.Equal(t, req2.Country, userAfterChange.Country.String)
	assert.Equal(t, req2.Nickname, userAfterChange.Nickname.String)
	assert.Equal(t, userBeforeChange.Password.String, userAfterChange.Password.String)
}

func TestModifyUserIntegration_ErrorModifyingUserThatDoeNotExist(t *testing.T) {
	client, conn, err := setupGRPCClient("localhost:9000")
	assert.NoError(t, err, "failed to set up gRPC client")
	defer conn.Close()

	d, err := setupPostgresDatasource()
	assert.NoError(t, err, "failed to connect to datasource")
	defer d.Disconnect()

	d.ResetUserStore()

	fakeId := "00000000-0000-0000-0000-000000000000"

	req2 := &api.ModifyUserRequest{
		Id:        fakeId,
		FirstName: "John",
		LastName:  "NoDoe",
		Email:     "john.nodoe@example.com",
		Country:   "UK",
		Nickname:  "johnnodoe",
	}

	resp2, err := client.ModifyUser(context.Background(), req2)
	assert.Error(t, err)
	assert.Nil(t, resp2)

	expectedError := fmt.Sprintf("User with id %s not found.", fakeId)
	assert.Contains(t, err.Error(), expectedError)
}

func TestDeleteUserIntegration_HappyPath(t *testing.T) {
	client, conn, err := setupGRPCClient("localhost:9000")
	assert.NoError(t, err, "failed to set up gRPC client")
	defer conn.Close()

	d, err := setupPostgresDatasource()
	assert.NoError(t, err, "failed to connect to datasource")
	defer d.Disconnect()

	d.ResetUserStore()

	req := &api.CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "password123",
		Country:   "US",
		Nickname:  "janedoe",
	}

	resp, err := client.CreateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Id, "response ID should not be empty")

	log.Print("Created user with ID: ", resp.Id)

	countBeforeDelete, err := d.GetUserCount()
	assert.NoError(t, err)

	user, err := d.GetUserById(resp.Id)
	assert.NoError(t, err)
	assert.Equal(t, resp.Id, user.ID.String)

	req2 := &api.DeleteUserRequest{
		Id: resp.Id,
	}

	resp2, err := client.DeleteUser(context.Background(), req2)
	assert.NoError(t, err)
	assert.Equal(t, "Successfully deleted user", resp2.Message)

	countAfterDelete, err := d.GetUserCount()
	assert.NoError(t, err)

	assert.Equal(t, countBeforeDelete-1, countAfterDelete)
}

func TestDeleteUserIntegration_ErrorDeletingUserThatDoeNotExist(t *testing.T) {
	client, conn, err := setupGRPCClient("localhost:9000")
	assert.NoError(t, err, "failed to set up gRPC client")
	defer conn.Close()

	d, err := setupPostgresDatasource()
	assert.NoError(t, err, "failed to connect to datasource")
	defer d.Disconnect()

	d.ResetUserStore()

	fakeId := "00000000-0000-0000-0000-000000000000"

	req2 := &api.DeleteUserRequest{
		Id: fakeId,
	}

	resp2, err := client.DeleteUser(context.Background(), req2)
	assert.Error(t, err)
	assert.Nil(t, resp2)

	expectedError := fmt.Sprintf("User with id %s not found.", fakeId)
	assert.Contains(t, err.Error(), expectedError)
}

func TestGetUsersIntegration_HappyPath(t *testing.T) {
	client, conn, err := setupGRPCClient("localhost:9000")
	assert.NoError(t, err, "failed to set up gRPC client")
	defer conn.Close()

	d, err := setupPostgresDatasource()
	assert.NoError(t, err, "failed to connect to datasource")
	defer d.Disconnect()

	d.ResetUserStore()

	req := &api.CreateUserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "password123",
		Country:   "US",
		Nickname:  "janedoe",
	}

	resp, err := client.CreateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Id, "response ID should not be empty")

	log.Print("Created user with ID: ", resp.Id)

	req2 := &api.CreateUserRequest{
		FirstName: "John",
		LastName:  "NoDoe",
		Email:     "john.nodoe@example.com",
		Password:  "password123",
		Country:   "UK",
		Nickname:  "johnnodoe",
	}

	resp2, err := client.CreateUser(context.Background(), req2)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp2.Id, "response ID should not be empty")

	log.Print("Created user with ID: ", resp2.Id)

	req3 := &api.CreateUserRequest{
		FirstName: "Joy",
		LastName:  "Doe",
		Email:     "joy.doe@example.com",
		Password:  "password123",
		Country:   "US",
		Nickname:  "joydoe",
	}

	resp3, err := client.CreateUser(context.Background(), req3)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp3.Id, "response ID should not be empty")

	log.Print("Created user with ID: ", resp3.Id)

	countOfUsers, err := d.GetUserCount()
	assert.NoError(t, err)
	assert.Equal(t, 3, countOfUsers)

	req4 := &api.GetUsersRequest{
		PageSize:      10,
		Page:          1,
		FilterCountry: "US",
	}

	resp4, err := client.GetUsers(context.Background(), req4)
	assert.NoError(t, err)
	assert.Len(t, resp4.Users, 2)
	// Order for function is done by created_at desc
	assert.Equal(t, resp4.Users[0].FirstName, "Joy")
	assert.Equal(t, resp4.Users[1].FirstName, "Jane")

	req5 := &api.GetUsersRequest{
		Page:     2,
		PageSize: 1,
	}

	resp5, err := client.GetUsers(context.Background(), req5)
	assert.NoError(t, err)
	assert.Len(t, resp5.Users, 1)
	assert.Equal(t, resp5.Users[0].FirstName, "John")

	req6 := &api.GetUsersRequest{
		Page:     3,
		PageSize: 1,
	}

	resp6, err := client.GetUsers(context.Background(), req6)
	assert.NoError(t, err)
	assert.Len(t, resp6.Users, 1)
	assert.Equal(t, resp6.Users[0].FirstName, "Jane")

	req7 := &api.GetUsersRequest{
		Page:           1,
		PageSize:       10,
		FilterLastName: "Doe",
	}

	resp7, err := client.GetUsers(context.Background(), req7)
	assert.NoError(t, err)
	// Fuzzy matching on last name is returning 3 users
	assert.Len(t, resp7.Users, 3)
}

func TestGetUsersIntegration_ErrorGettingUsersWithInvalidPage(t *testing.T) {
	client, conn, err := setupGRPCClient("localhost:9000")
	assert.NoError(t, err, "failed to set up gRPC client")
	defer conn.Close()

	req := &api.GetUsersRequest{
		Page:     -1,
		PageSize: 1,
	}

	resp, err := client.GetUsers(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Contains(t, err.Error(), "page must be >= 1")
}

func TestGetUsersIntegration_ErrorGettingUsersWithInvalidPageSize(t *testing.T) {
	client, conn, err := setupGRPCClient("localhost:9000")
	assert.NoError(t, err, "failed to set up gRPC client")
	defer conn.Close()

	req := &api.GetUsersRequest{
		Page:     1,
		PageSize: -1,
	}

	resp, err := client.GetUsers(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Contains(t, err.Error(), "page_size must be >= 1")
}

func TestGetUsersIntegration_ErrorGettingUsersWithNoMatchingFilter(t *testing.T) {
	client, conn, err := setupGRPCClient("localhost:9000")
	assert.NoError(t, err, "failed to set up gRPC client")
	defer conn.Close()

	req := &api.GetUsersRequest{
		Page:            1,
		PageSize:        1,
		FilterFirstName: "ZZ",
		FilterId:        "00000000-0000-0000-0000-000000000000",
		FilterLastName:  "ZZ",
		FilterCountry:   "ZZ",
		FilterEmail:     "ZZ",
		FilterNickname:  "ZZ",
	}

	resp, err := client.GetUsers(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	assert.Contains(t, err.Error(), "no users found for supplied filters")
}
