package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set up test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("SERVER_ADDR", "127.0.0.1:8080")
	os.Setenv("BASE_URL", "http://example.com")
	os.Setenv("PATH_PREFIX", "/api/v1")

	// Load the configuration
	InitConfig()

	assert.Equal(t, "localhost", Envs.DBHost)
	assert.Equal(t, "5432", Envs.DBPort)
	assert.Equal(t, "testuser", Envs.DBUser)
	assert.Equal(t, "testpass", Envs.DBPassword)
	assert.Equal(t, "testdb", Envs.DBName)
	assert.Equal(t, "127.0.0.1:8080", Envs.ServerAddress)
	assert.Equal(t, "http://example.com", Envs.BaseUrl)
	assert.Equal(t, "/api/v1", Envs.PathPrefix)
}

func TestGetDataSourceName(t *testing.T) {
	// Set up test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")

	// Load the configuration
	InitConfig()

	expectedDSN := "testuser:testpass@tcp(localhost:5432)/testdb?checkConnLiveness=false&parseTime=true"
	assert.Equal(t, expectedDSN, GetDataSourceName())
}
