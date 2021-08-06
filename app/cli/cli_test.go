package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/razielsd/antibruteforce/app/config"
)

func TestNewCli(t *testing.T) {
	cfg := config.AppConfig{Addr: "0.0.0.0:8081"}
	appCli := NewCli(cfg)
	require.NotNil(t, appCli)
}
