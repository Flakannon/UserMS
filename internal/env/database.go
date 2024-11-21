package env

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

// DatabaseConfig holds the database configuration.
type DatabaseConfig struct {
	Host     string `mapstructure:"HOST"`
	Username string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	Port     string `mapstructure:"PORT"`
	Database string `mapstructure:"DATABASE"`
	Schema   string `mapstructure:"SCHEMA"`
}

// Validate ensures all required fields in the DatabaseConfig are set.
func (c DatabaseConfig) Validate() error {
	if c.Host == "" || c.Username == "" || c.Password == "" || c.Port == "" || c.Database == "" || c.Schema == "" {
		return fmt.Errorf("one or more required fields are missing")
	}
	return nil
}

func LoadDatabaseConfig() (config DatabaseConfig, err error) {
	if err = viperBindDB("POSTGRES", &config); err != nil {
		return DatabaseConfig{}, fmt.Errorf("failed to load database configs for prefix %s: %w", "POSTGRES", err)
	}
	if err = config.Validate(); err != nil {
		return DatabaseConfig{}, fmt.Errorf("validation failed  %s: %w", "POSTGRES", err)
	}

	slog.Info("Loaded configuration for prefix %s: %+v\n", "POSTGRES", config)

	return
}

func viperBindDB(prefix string, config *DatabaseConfig) error {
	viper.SetEnvPrefix(prefix)
	viper.AutomaticEnv()

	viper.BindEnv("HOST")
	viper.BindEnv("USER")
	viper.BindEnv("PASSWORD")
	viper.BindEnv("PORT")
	viper.BindEnv("DATABASE")
	viper.BindEnv("WAREHOUSE")
	viper.BindEnv("SCHEMA")
	viper.BindEnv("ROLE")

	return viper.Unmarshal(&config)
}
