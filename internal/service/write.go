package service

import (
	"context"

	"github.com/EFG/internal/datasource/dto"
)

type Writer interface {
	CreateUser(ctx context.Context, user dto.UserDTO) (string, error)
	ModifyUser(ctx context.Context, user dto.UserDTO) error
	DeleteUser(ctx context.Context, userID int) error
}

func FormatUserAndPersist(ctx context.Context, writer Writer, user User) (id string, err error) {
	user.hashPassword()

	userEntityToWrite := user.toDTO()

	id, err = writer.CreateUser(ctx, userEntityToWrite)
	if err != nil {
		return "", err
	}

	return
}
