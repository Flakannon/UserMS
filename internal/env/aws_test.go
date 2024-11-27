package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAWSConfig(t *testing.T) {
	os.Setenv("AWS_USER_CHANGE_NOTIFICATION_TOPIC", "topic-name")
	os.Setenv("AWS_LOCALSTACK_URL", "localhost")
	os.Setenv("AWS_REGION", "us-west-2")

	config, err := LoadAWSConfig()
	assert.NoError(t, err)

	assert.Equal(t, "topic-name", config.UserChangeNotificationTopic)
	assert.Equal(t, "localhost", config.LocalstackURL)
	assert.Equal(t, "us-west-2", config.Region)
}
