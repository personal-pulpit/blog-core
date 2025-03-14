package logger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetZapLoggerInstance(t *testing.T) {
	logger := GetZapLoggerInstance()
	logger.Info("Hi Its a test", "test time", time.Now().String())
	require.NotNil(t, logger)
}
