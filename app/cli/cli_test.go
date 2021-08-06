package cli

import (
	"github.com/razielsd/antibruteforce/app/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewCli(t *testing.T) {
	cfg := config.AppConfig{Addr: "0.0.0.0:8080"}
	appCli := NewCli(cfg)
	require.NotNil(t, appCli)
}
