package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

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
		if strings.Contains(err.Error(), "user_email_unique") {
			return "", fmt.Errorf("email already exists: %s", user.Email.String)
		}
		return "", fmt.Errorf("database error: %w", err)
	}

	return id, nil
}

//go:embed scripts/postgres_update_user_function_call.sql
var updateUserFunctionCall string

func (d *Client) ModifyUser(ctx context.Context, user dto.UserDTO) error {
	slog.Info("Modifying user", "id", user.ID)
	_, err := d.DB.ExecContext(ctx, updateUserFunctionCall,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Nickname,
		user.Password,
		user.Email,
		user.Country,
	)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}

//go:embed scripts/postgres_delete_user_function_call.sql
var deleteUserFunctionCall string

func (d *Client) DeleteUser(ctx context.Context, userUUID string) error {
	_, err := d.DB.ExecContext(ctx, deleteUserFunctionCall, userUUID)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}
