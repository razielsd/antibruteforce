package config

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	testCfg := AppConfig{
		Port:          8090,
		RateLogin:     12,
		RatePwd:       121,
		RateIP:        10001,
		Whitelist:     []string{"192.168.1.1", "192.168.1.72"},
		Blacklist:     []string{"10.10.1.1", "10.10.1.72"},
		WhitelistPath: "/tmp/whitelist.txt",
		BlacklistPath: "/tmp/blacklist.txt",
	}
	os.Setenv("ABF_PORT", strconv.Itoa(testCfg.Port))
	os.Setenv("ABF_RATE_LOGIN", strconv.Itoa(testCfg.RateLogin))
	os.Setenv("ABF_RATE_PWD", strconv.Itoa(testCfg.RatePwd))
	os.Setenv("ABF_RATE_IP", strconv.Itoa(testCfg.RateIP))
	os.Setenv("ABF_WHITELIST", strings.Join(testCfg.Whitelist, ","))
	os.Setenv("ABF_BLACKLIST", strings.Join(testCfg.Blacklist, ","))
	os.Setenv("ABF_WHITELIST_PATH", testCfg.WhitelistPath)
	os.Setenv("ABF_BLACKLIST_PATH", testCfg.BlacklistPath)
	cfg, err := GetConfig()
	require.NoError(t, err)

	require.Equal(t, testCfg, cfg)
}
