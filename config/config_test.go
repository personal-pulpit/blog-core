package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDevelopmentConfig(t *testing.T) {
	config := GetConfigInstance()

	require.NotEmpty(t, config)
	require.Equal(t, GetEnv(), Development)
}
