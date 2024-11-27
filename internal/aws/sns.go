package aws

import (
	"context"
	"fmt"

	"github.com/EFG/internal/env"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	snspkg "github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNS struct {
	client *snspkg.Client
	Config env.AWSConfig
	ctx    context.Context
}

func snsClientWithResolution(cfg aws.Config, localStackURL string) *snspkg.Client {
	return snspkg.NewFromConfig(cfg, func(o *snspkg.Options) {
		if localStackURL != "" {
			o.BaseEndpoint = aws.String(localStackURL)
		}
	})
}

func NewSNSClient(ctx context.Context, cfg env.AWSConfig) (sns SNS, err error) {
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.Region))
	if err != nil {
		return sns, fmt.Errorf("issue loading AWS SDK config: %w", err)
	}

	snsClient := snsClientWithResolution(sdkConfig, cfg.LocalstackURL)

	sns = SNS{
		client: snsClient,
		Config: cfg,
		ctx:    ctx,
	}

	return sns, nil
}

func (s *SNS) PublishMessage(message []byte, topicARN string) (string, error) {
	output, err := s.client.Publish(s.ctx, &snspkg.PublishInput{
		TopicArn: aws.String(topicARN),
		Message:  aws.String(string(message)),
	})
	if err != nil {
		return "", fmt.Errorf("error publishing message to SNS: %w", err)
	}

	return *output.MessageId, nil
}

func (s *SNS) SubscribeToTopic(topicARN, protocol, endpoint string) (string, error) {
	output, err := s.client.Subscribe(s.ctx, &snspkg.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String(protocol),
		Endpoint: aws.String(endpoint),
	})
	if err != nil {
		return "", fmt.Errorf("error subscribing to topic: %w", err)
	}

	return *output.SubscriptionArn, nil
}
