package env

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

// DatabaseConfig holds a database configuration.
type AWSConfig struct {
	UserChangeNotificationTopic string `mapstructure:"USER_CHANGE_NOTIFICATION_TOPIC"`
	LocalstackURL               string `mapstructure:"LOCALSTACK_URL"`
	Region                      string `mapstructure:"REGION"`
}

func LoadAWSConfig() (config AWSConfig, err error) {
	if err = viperBindAWS("AWS", &config); err != nil {
		return AWSConfig{}, fmt.Errorf("failed to load AWS configs for prefix %s: %w", "AWS", err)
	}

	slog.Info("Loaded AWS configuration",
		"prefix", "AWS",
		"topic", config.UserChangeNotificationTopic,
		"localstackURL", config.LocalstackURL,
		"region", config.Region)

	return
}

func viperBindAWS(prefix string, config *AWSConfig) error {
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()

	viper.BindEnv("USER_CHANGE_NOTIFICATION_TOPIC")
	viper.BindEnv("LOCALSTACK_URL")
	viper.BindEnv("REGION")

	return viper.Unmarshal(&config)
}

func (a *AWSConfig) IsValid() bool {
	return a.UserChangeNotificationTopic != "" && a.LocalstackURL != "" && a.Region != ""
}
