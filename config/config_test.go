package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDevelopmentConfig(t *testing.T) {
	os.Setenv("ENV","development")

	config := GetConfigInstance()
	
	require.NotEmpty(t, config)
	require.Equal(t, GetEnv(), Development)
}
