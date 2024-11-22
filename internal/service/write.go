package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/EFG/internal/datasource/dto"
)

type Writer interface {
	CreateUser(ctx context.Context, user dto.UserDTO) (string, error)
	ModifyUser(ctx context.Context, user dto.UserDTO) error
	DeleteUser(ctx context.Context, userUUID string) error
}

func FormatNewUserAndPersist(ctx context.Context, writer Writer, user User) (id string, err error) {
	err = user.hashPassword()
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	userEntityToWrite := user.toDTO()

	id, err = writer.CreateUser(ctx, userEntityToWrite)
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return
}

func FormatExistingUserAndPersist(ctx context.Context, writer Writer, user User) error {
	// if password is part of the modification, hash it
	if user.Password != "" {
		err := user.hashPassword()
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
	}

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
