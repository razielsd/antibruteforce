// +build e2e

package e2e

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCli(t *testing.T) {
	bin, srvCmd, srvCancel := runAbf(t)
	defer srvCancel()
	err := srvCmd.Start()
	require.NoError(t, err)

	tests := []struct {
		name string
		text string
		arg  []string
	}{
		{
			name: "no args",
			text: "default",
			arg:  []string{""},
		},
		{
			name: "show whitelist no args",
			text: "whitelist",
			arg:  []string{"whitelist"},
		},
		{
			name: "show blacklist no args",
			text: "blacklist",
			arg:  []string{"blacklist"},
		},
		{
			name: "show empty whitelist",
			text: "whitelist_show_empty",
			arg:  []string{"whitelist", "show"},
		},
		{
			name: "show empty blacklist",
			text: "blacklist_show_empty",
			arg:  []string{"blacklist", "show"},
		},
		{
			name: "whitelist add ip",
			text: "ok",
			arg:  []string{"whitelist", "add", "10.10.1.12"},
		},
		{
			name: "show whitelist with ip",
			text: "whitelist_with_ip",
			arg:  []string{"whitelist", "show"},
		},
		{
			name: "whitelist remove ip",
			text: "ok",
			arg:  []string{"whitelist", "rm", "10.10.1.12"},
		},
		{
			name: "show empty whitelist after remove ip",
			text: "whitelist_show_empty",
			arg:  []string{"whitelist", "show"},
		},
		{
			name: "blacklist add ip",
			text: "ok",
			arg:  []string{"blacklist", "add", "192.168.1.10"},
		},
		{
			name: "show blacklist with ip",
			text: "blacklist_with_ip",
			arg:  []string{"blacklist", "show"},
		},
		{
			name: "blacklist remove ip",
			text: "ok",
			arg:  []string{"blacklist", "rm", "192.168.1.10"},
		},
		{
			name: "show empty blacklist after remove ip",
			text: "blacklist_show_empty",
			arg:  []string{"blacklist", "show"},
		},
		{
			name: "drop bucket by login",
			text: "ok",
			arg:  []string{"bucket", "drop", "login", "Ivan"},
		},
		{
			name: "drop bucket by pwd",
			text: "ok",
			arg:  []string{"bucket", "drop", "pwd", "123456"},
		},
		{
			name: "drop bucket by ip",
			text: "ok",
			arg:  []string{"bucket", "drop", "ip", "192.168.1.10"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, bin, test.arg...)
			cmd.Wait()

			out, err := cmd.Output()
			require.NoError(t, err)
			require.Equal(t, getExpectedOutput(test.text), string(out))
		})
	}
}

func getExpectedOutput(name string) string {
	output := make(map[string]string)
	output["default"] = "Antibruteforce service cli\n\nUsage:\n  abf [command]\n\nAvailable Commands:\n  blacklist   Show/add/remove blacklist\n  bucket      Drop bucket by login, password or ip\n  completion  generate the autocompletion script for the specified shell\n  help        Help about any command\n  server      Run service\n  version     Show version\n  whitelist   Show/add/remove whitelist\n\nFlags:\n  -h, --help   help for abf\n\nUse \"abf [command] --help\" for more information about a command.\n"
	output["whitelist"] = "Show/add/remove whitelist\n\nUsage:\n  abf whitelist [command]\n\nAvailable Commands:\n  add         Add ip/subnet to whitelist\n  rm          Remove ip/subnet from whitelist\n  show        Show whitelist\n\nFlags:\n  -h, --help   help for whitelist\n\nUse \"abf whitelist [command] --help\" for more information about a command.\n"
	output["blacklist"] = "Show/add/remove blacklist\n\nUsage:\n  abf blacklist [command]\n\nAvailable Commands:\n  add         Add ip/subnet to blacklist\n  rm          Remove ip/subnet from blacklist\n  show        Show blacklist\n\nFlags:\n  -h, --help   help for blacklist\n\nUse \"abf blacklist [command] --help\" for more information about a command.\n"
	output["whitelist_show_empty"] = "--=== Whitelist ===--\nEmpty\n\n"
	output["blacklist_show_empty"] = "--=== Blacklist ===--\nEmpty\n\n"
	output["ok"] = "OK\n"
	output["blacklist_with_ip"] = "--=== Blacklist ===--\n192.168.1.10\n\n"
	output["whitelist_with_ip"] = "--=== Whitelist ===--\n10.10.1.12\n\n"
	return output[name]
}

func runAbf(t *testing.T) (string, *exec.Cmd, context.CancelFunc) {
	bin := os.Getenv("ABF_BIN")
	require.NotEmpty(t, bin, "Empty env ABF_BIN")
	if _, err := os.Stat(bin); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("executable not found in env ABF_BIN: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, bin, "server")
	cmd.Env = os.Environ()
	return bin, cmd, cancel
}
