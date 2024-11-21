package postgres

import (
	"context"

	"github.com/EFG/internal/datasource/dto"
)

func (d *Client) GetUsers(ctx context.Context, user dto.GetUsersArgs) (dto.UsersDTO, int, error) {
	return nil, 0, nil
}
