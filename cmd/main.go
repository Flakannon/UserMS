package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/EFG/api"
	"github.com/EFG/internal/datasource/database/postgres"
	"github.com/EFG/internal/env"
	"github.com/EFG/internal/logger"
	"github.com/EFG/internal/server"
	"google.golang.org/grpc"
)

func main() {
	slog.SetDefault(logger.SetUpLogger(logger.LoggerInitOpts{
		Writer:         os.Stdout,
		VerbosityLevel: 0,
	}))

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to listen on port 9000: %w", err))
	}

	grpcServer := grpc.NewServer()

	postgresConfig, err := env.LoadDatabaseConfig()
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to load database config: %w", err))
	}
	postgresDataSource := postgres.NewClient(postgresConfig)
	err = postgresDataSource.Connect()
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}
	defer postgresDataSource.Close()

	userServer := server.NewServer(postgresDataSource)

	api.RegisterUserServiceServer(grpcServer, userServer)

	slog.Info("gRPC server is listening on port 9000")
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal(fmt.Errorf("failed to serve: %w", err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	grpcServer.GracefulStop()
	slog.Info("gRPC server is shutting down")
}

// user management micro service - implement in go

// Requirements

// A user will be stored using the following schema: - done hashing needs to be implemented at app level
// ID: a unique identifier for the user
// FirstName: the user's first name
// LastName: the user's last name
// Nickname: the user's nickname
// Password[Hashed]: the user's password, hashed using bcrypt
// Email: the user's email address
// Country: the user's country
// CreatedAt: the date and time the user was created
// UpdatedAt: the date and time the user was last updated

// Return a paginated list of users, allowing the results to be filtered by any of the fields e.g Country = "UK" - db done

// The service must allow the following operations:
// Create a new user - db done
// Modify an existing user - db done
// Remove a user - db done
// Return a paginated list of users, allowing the results to be filtered by any of the fields e.g Country = "UK" - db done

// The service must
// Provide an HTTP or gRPC API to interact with the service
// Use a sensible storage mechanism for storing users - done
// Have the ability to notify other services of changes to user entities
// Have meaningful logs
// Be well documented - a good balance between code comments and external documentation in readme (especially of choices made and why)
// Have a health check

// The service must NOT:
// Provide authentication or authorisation
