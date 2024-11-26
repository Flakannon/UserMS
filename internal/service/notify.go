package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
)

type Notifier interface {
	PublishUserChange(ctx context.Context, message []byte) error
}

func NotifyOfUserChange(ctx context.Context, notifier Notifier, changeData UserChange) error {
	slog.Info("Notifying of user change", "changeData", changeData)

	notificationMessage, err := json.Marshal(changeData)
	if err != nil {
		slog.Error("failed to marshal user change data", "error", err)
		return fmt.Errorf("failed to marshal user change data: %w", err)
	}

	if err := notifier.PublishUserChange(ctx, notificationMessage); err != nil {
		slog.Error("failed to publish user change notification", "message", notificationMessage, "error", err)
		return fmt.Errorf("failed to publish user change notification but have placed in backup for team to review: %w", err)

	}

	slog.Info("User change notification complete")
	return nil
}
