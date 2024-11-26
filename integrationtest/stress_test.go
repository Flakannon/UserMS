package integrationtest

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/EFG/api"
	"github.com/stretchr/testify/assert"
)

func TestStressForCreateUsers(t *testing.T) {
	concurrentClients := 15
	totalRequests := 1000

	client, conn, err := setupGRPCClient("localhost:9000")
	assert.NoError(t, err, "failed to set up gRPC client")
	defer conn.Close()

	d, err := setupPostgresDatasource()
	assert.NoError(t, err, "failed to connect to datasource")
	defer d.Disconnect()

	err = d.ResetUserStore()
	assert.NoError(t, err)

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(concurrentClients)

	resultsChan := make(chan time.Duration, totalRequests)

	emailIter := 0

	for i := 0; i < concurrentClients; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < totalRequests/concurrentClients; j++ {
				start := time.Now()

				mu.Lock()
				currentEmail := emailIter
				emailIter++
				mu.Unlock()

				_, err := client.CreateUser(context.Background(), &api.CreateUserRequest{
					FirstName: "Test",
					LastName:  "User",
					Email:     fmt.Sprintf("testuser%d@example.com", currentEmail),
					Password:  "password123",
					Country:   "US",
					Nickname:  "testuser",
				})
				duration := time.Since(start)
				resultsChan <- duration

				if err != nil {
					log.Printf("Request failed: %v", err)
				}

				// throttle requests to avoid overwhelming the server
				time.Sleep(250 * time.Millisecond)
			}
		}()
	}

	wg.Wait()
	close(resultsChan)

	var numRequests int

	for range resultsChan {
		numRequests++
	}

	count, err := d.GetUserCount()
	assert.NoError(t, err, "failed to count users")
	assert.Equal(t, numRequests, count, "number of users created does not match expected count")
}
