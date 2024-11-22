package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EFG/internal/datasource/database"
	"github.com/EFG/internal/datasource/dto"
	"github.com/EFG/internal/utils"
)

type MockClient struct {
	database.BaseClient
	UUID                    string
	WriteUserRowCallHistory dto.UsersDTO
	TestRequiresError       bool
}

func (m *MockClient) CreateUser(ctx context.Context, user dto.UserDTO) (string, error) {
	if m.TestRequiresError {
		return "", fmt.Errorf("mock db error for create user")
	}
	m.UserWritten(user)
	return m.UUID, nil
}

func (m *MockClient) ModifyUser(ctx context.Context, user dto.UserDTO) error {
	if m.TestRequiresError {
		return fmt.Errorf("mock db error for modify user")
	}
	m.UserWritten(user)

	return nil
}

func (m *MockClient) DeleteUser(ctx context.Context, id string) error {
	if m.TestRequiresError {
		return fmt.Errorf("mock db error for delete user")
	}
	m.UserWritten(dto.UserDTO{ID: utils.ToNullString(id)})
	return nil
}

func (d *MockClient) UserWritten(user dto.UserDTO) {
	d.WriteUserRowCallHistory = append(d.WriteUserRowCallHistory, user)
}

func (d *MockClient) Reset() {
	d.WriteUserRowCallHistory = dto.UsersDTO{}
}

func (d *MockClient) GetCallHistory() dto.UsersDTO {
	return d.WriteUserRowCallHistory
}

func mockSQLRowsGetUsersFromDataSource() *sql.Rows {
	timestamp := time.Now()
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "email", "country", "created_at", "updated_at"}).
		AddRow("1", "John", "Doe", "johndoe", "john.doe@example.com", "US", timestamp, timestamp).
		AddRow("2", "Jane", "Doe", "janedoe", "jane.doe@example.com", "US", timestamp, timestamp)
	return database.MockRowsToSQLRows(rows)
}

func (m *MockClient) GetUsers(ctx context.Context, user dto.GetUsersArgs) (dto.UsersDTO, int, error) {
	if m.TestRequiresError {
		return nil, 0, fmt.Errorf("mock db error for get users")
	}
	rows := mockSQLRowsGetUsersFromDataSource()

	users, err := scanUsers(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to scan users: %w", err)
	}

	return users, len(users), nil
}
