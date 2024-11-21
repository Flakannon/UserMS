package service

import (
	"context"

	"github.com/EFG/internal/datasource/dto"
)

type Reader interface {
	GetUsers(ctx context.Context, args dto.GetUsersArgs) (dto.UsersDTO, int, error)
}

func GetPaginatedUsersList(ctx context.Context, reader Reader, args dto.GetUsersArgs) (dto.UsersDTO, int, error) {
	users, total, err := reader.GetUsers(ctx, args)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
