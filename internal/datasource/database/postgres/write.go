package postgres

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/EFG/internal/datasource/dto"
)

//go:embed scripts/postgres_create_user_function_call.sql
var createUserFunctionCall string

func (d *Client) CreateUser(ctx context.Context, user dto.UserDTO) (string, error) {
	var id string

	err := d.DB.QueryRowContext(ctx, createUserFunctionCall,
		user.FirstName,
		user.LastName,
		user.Nickname,
		user.Password,
		user.Email,
		user.Country,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (d *Client) ModifyUser(ctx context.Context, user dto.UserDTO) error {
	return nil
}

func (d *Client) DeleteUser(ctx context.Context, userID int) error {
	return nil
}
