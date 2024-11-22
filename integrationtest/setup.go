package integrationtest

import (
	"fmt"

	"github.com/EFG/api"
	"github.com/EFG/internal/datasource/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Datasource interface {
	GetUserById(string) (dto.UserDTO, error)
	GetUserCount() (int, error)
	ResetUserStore() error
	Disconnect() error
}

type Notifier interface {
	GetNotifications() ([]string, error)
}

func setupPostgresDatasource() (Datasource, error) {
	client := &PostgresClient{}
	err := client.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to datasource: %w", err)
	}
	return client, nil
}

// create grpc client server to connect to grpc server in docker
func setupGRPCClient(address string) (api.UserServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to gRPC server at %s: %w", address, err)
	}

	client := api.NewUserServiceClient(conn)

	return client, conn, nil
}
