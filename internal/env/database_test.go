package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDatabaseConfigs(t *testing.T) {
	// Set environment variables for POSTGRES
	os.Setenv("POSTGRES_HOST", "localhost")
	os.Setenv("POSTGRES_USER", "user")
	os.Setenv("POSTGRES_PASSWORD", "password")
	os.Setenv("POSTGRES_PORT", "5432")
	os.Setenv("POSTGRES_DATABASE", "testdb")
	os.Setenv("POSTGRES_SCHEMA", "public")

	config, err := LoadDatabaseConfigs()
	assert.NoError(t, err)
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, "user", config.Username)
	assert.Equal(t, "password", config.Password)
	assert.Equal(t, "5432", config.Port)
	assert.Equal(t, "testdb", config.Database)
	assert.Equal(t, "public", config.Schema)

	// Unset environment variables
	os.Unsetenv("POSTGRES_HOST")
	os.Unsetenv("POSTGRES_USER")
	os.Unsetenv("POSTGRES_PASSWORD")
	os.Unsetenv("POSTGRES_PORT")
	os.Unsetenv("POSTGRES_DATABASE")
	os.Unsetenv("POSTGRES_SCHEMA")
}
