package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/EFG/internal/datasource/dto"
)

type Reader interface {
	GetUsers(ctx context.Context, args dto.GetUsersArgs) (dto.UsersDTO, int, error)
}

func GetPaginatedUsersList(ctx context.Context, reader Reader, args dto.GetUsersArgs) (dto.UsersDTO, int, error) {
	users, total, err := reader.GetUsers(ctx, args)
	if err != nil {
		slog.Error("failed to get users", "error", err)
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}
