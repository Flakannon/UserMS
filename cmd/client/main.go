package main

import (
	"context"
	"log"
	"time"

	"github.com/EFG/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
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
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	log.Printf("Created user: %s", createResp.Id)

	// Delete the user
	deleteResp, err := client.DeleteUser(ctx, &api.DeleteUserRequest{Id: createResp.Id})
	if err != nil {
		log.Fatalf("DeleteUser failed: %v", err)
	}
	log.Printf("Deleted user: %s", deleteResp.Message)
}
