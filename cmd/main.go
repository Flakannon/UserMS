package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/EFG/api"
	"github.com/EFG/internal/aws"
	"github.com/EFG/internal/datasource/database/postgres"
	"github.com/EFG/internal/env"
	"github.com/EFG/internal/logger"
	"github.com/EFG/internal/notifier"
	"github.com/EFG/internal/server"
	"github.com/EFG/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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

	awsConfig, err := env.LoadAWSConfig()
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to load AWS config: %w", err))
	}

	ctx := context.Background()

	// Always default to NoOpNotifier
	var notifierService service.Notifier = notifier.NewNoOpNotifier()

	if awsConfig.IsValid() {
		snsClient, err := aws.NewSNSClient(ctx, awsConfig)
		if err != nil {
			slog.Info("Failed to create SNS client, using NoOpNotifier", "error", err)
		} else {
			notifierService = notifier.NewSNSNotifier(snsClient)
			slog.Info("SNS notifier service is enabled")
		}
	}

	userServer := server.NewServer(postgresDataSource, notifierService, time.Now)

	api.RegisterUserServiceServer(grpcServer, userServer)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	// check heartbeat for database to update health status
	go func() {
		for {
			slog.Info("Health check for critical connections running")
			err := postgresDataSource.PingDatabase()
			if err != nil {
				slog.Warn("Database health check failed", "error", err)
				healthServer.SetServingStatus("api.UserService", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			} else {
				slog.Info("Database is healthy")
				healthServer.SetServingStatus("api.UserService", grpc_health_v1.HealthCheckResponse_SERVING)
			}
			time.Sleep(10 * time.Second)
		}
	}()

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
