package logger

import (
	"blog/config"
	"testing"
	"github.com/stretchr/testify/require"

)

func TestGetZapLoggerInstance(t *testing.T){
	logger := GetZapLoggerInstance(&config.GetConfigInstance().Logger)

	require.NotNil(t,logger)
}