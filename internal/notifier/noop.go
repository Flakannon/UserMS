package notifier

import (
	"context"
	"log/slog"
)

type NoOpNotifier struct{}

func NewNoOpNotifier() *NoOpNotifier {
	return &NoOpNotifier{}
}

func (n *NoOpNotifier) PublishUserChange(ctx context.Context, message []byte) error {
	slog.Info("NoOpNotifier: would have published", "message", string(message))
	return nil
}
