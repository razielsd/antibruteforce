package logger

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/razielsd/antibruteforce/app/config"
)

func TestLogger_ValidLogLevel_GetLogger(t *testing.T) {
	cfg := config.AppConfig{
		LogLevel: "WARN",
	}
	logger, err := GetLogger(cfg)
	require.NoError(t, err)
	require.NotEmpty(t, logger)
}

func TestLogger_InvalidLogLevel_GetError(t *testing.T) {
	cfg := config.AppConfig{
		LogLevel: "WARNing",
	}
	_, err := GetLogger(cfg)
	require.Error(t, err)
}
