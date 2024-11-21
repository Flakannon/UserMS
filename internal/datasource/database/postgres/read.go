package postgres

import (
	"context"

	"github.com/EFG/internal/datasource/dto"
)

func (d *Client) GetUsers(ctx context.Context, user dto.UserDTO) error {
	return nil
}
