package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/EFG/api"
	"github.com/EFG/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to connect: %w", err))
	}
	client := api.NewUserServiceClient(conn)

	// Create a new user
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	createResp, err := client.CreateUser(ctx, &api.CreateUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "securepassword",
		Country:   "US",
		Nickname:  "johndoe",
	})
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create user: %w", err))
	}

	slog.Info("Created user", "id", createResp.Id)

	// Delete the user
	deleteResp, err := client.DeleteUser(ctx, &api.DeleteUserRequest{Id: createResp.Id})
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to delete user: %w", err))
	}
	slog.Info("Deleted user", "id", deleteResp.Message)
}
