package postgres

import (
	"context"

	"github.com/EFG/internal/datasource/dto"

	_ "embed"
)

// p_ID UUID DEFAULT NULL,
// p_country TEXT DEFAULT NULL,
// p_email TEXT DEFAULT NULL,
// p_first_name TEXT DEFAULT NULL,
// p_last_name TEXT DEFAULT NULL,
// p_nick_name TEXT DEFAULT NULL,
// p_page INT DEFAULT 1,
// p_page_size INT DEFAULT 10

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
		return nil, 0, err
	}
	defer rows.Close()

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
			return nil, 0, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, len(users), nil
}
