package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/EFG/internal/datasource/dto"

	_ "embed"
)

//go:embed scripts/postgres_get_users_function_call.sql
var getUsersFunctionCall string

func (d *Client) GetUsers(ctx context.Context, user dto.GetUsersArgs) (dto.UsersDTO, int, error) {
	rows, err := d.DB.QueryContext(ctx, getUsersFunctionCall,
		user.FilterID,
		user.FilterCountry,
		user.FilterEmail,
		user.FilterFirstName,
		user.FilterLastName,
		user.FilterNickname,
		user.Page,
		user.PageSize,
	)
	if err != nil {
		slog.Error("failed to call get_users function", "error", err)
		return nil, 0, fmt.Errorf("failed to call get_users function: %w", err)
	}
	defer rows.Close()

	users, err := scanUsers(rows)
	if err != nil {
		slog.Error("failed to scan users", "error", err)
		return nil, 0, fmt.Errorf("failed to scan users: %w", err)
	}

	if len(users) == 0 {
		return nil, 0, fmt.Errorf("no users found for supplied filters")
	}

	return users, len(users), nil
}

func scanUsers(rows *sql.Rows) (dto.UsersDTO, error) {
	var users dto.UsersDTO
	for rows.Next() {
		var u dto.UserDTO
		if err := rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Nickname,
			&u.Email,
			&u.Country,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
