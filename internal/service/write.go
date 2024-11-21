package service

import (
	"context"
	"fmt"

	"github.com/EFG/internal/datasource/dto"
)

type Writer interface {
	CreateUser(ctx context.Context, user dto.UserDTO) (string, error)
	ModifyUser(ctx context.Context, user dto.UserDTO) error
	DeleteUser(ctx context.Context, userUUID string) error
}

func FormatNewUserAndPersist(ctx context.Context, writer Writer, user User) (id string, err error) {
	user.hashPassword()

	userEntityToWrite := user.toDTO()

	id, err = writer.CreateUser(ctx, userEntityToWrite)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return
}

func FormatExistingUserAndPersist(ctx context.Context, writer Writer, user User) error {
	userEntityToWrite := user.toDTO()

	err := writer.ModifyUser(ctx, userEntityToWrite)
	if err != nil {
		return fmt.Errorf("failed to modify user: %w", err)
	}

	return nil
}

func DeleteUserFromDatasource(ctx context.Context, writer Writer, userUUID string) error {
	err := writer.DeleteUser(ctx, userUUID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
