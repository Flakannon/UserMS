package notifier

import (
	"context"
	"log/slog"

	"github.com/EFG/internal/aws"
)

type SNSNotifier struct {
	snsClient aws.SNS
}

func NewSNSNotifier(notification aws.SNS) *SNSNotifier {
	return &SNSNotifier{
		snsClient: notification,
	}
}

// PublishUserChange sends the notification via SNS
func (n *SNSNotifier) PublishUserChange(ctx context.Context, message []byte) error {
	messageID, err := n.snsClient.PublishMessage(message, n.snsClient.Config.UserChangeNotificationTopic)
	if err != nil {
		return err
	}

	slog.Info("Published message to SNS", "messageID", messageID)
	return nil
}
