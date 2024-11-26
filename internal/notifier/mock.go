package notifier

import (
	"context"
	"errors"
)

type MockNotifier struct {
	PublishCalled            bool
	PublishedMessages        [][]byte
	TestRequiresPublishError bool
	TestRequiresDLQError     bool
}

func (m *MockNotifier) PublishUserChange(ctx context.Context, message []byte) error {
	m.PublishCalled = true

	if m.TestRequiresPublishError {
		return errors.New("mock publish error")
	}

	m.PublishedMessages = append(m.PublishedMessages, message)

	return nil
}
