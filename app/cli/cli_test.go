package cli

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/razielsd/antibruteforce/app/api"
	"github.com/razielsd/antibruteforce/app/config"
)

func TestNewCli(t *testing.T) {
	cfg := config.AppConfig{Addr: "0.0.0.0:8081"}
	appCli := NewCli(cfg)
	require.NotNil(t, appCli)
}

func TestCli_ShowWhitelist(t *testing.T) {
	cfg := config.AppConfig{Addr: "0.0.0.0:8081"}
	ips := []string{"192.168.1.71", "10.10.1.15"}
	expResp := api.SuccessResponse{
		Result: ips,
	}
	expTxt := "--=== Whitelist ===--\n" + strings.Join(ips, "\n") + "\n"
	mock := newAPIClientMock(http.StatusOK, expResp.JSON(), nil)
	appCli := NewCli(cfg)
	appCli.client.setHTTPClient(mock)
	require.NotNil(t, appCli)
	txt, err := appCli.ShowWhitelist()
	require.NoError(t, err)
	require.Equal(t, expTxt, txt)
}

func TestCli_Success(t *testing.T) {
	ips := []string{"192.168.1.71", "10.10.1.15"}
	ip := "172.168.10.5"
	cfg := config.AppConfig{Addr: "0.0.0.0:8081"}
	appCli := NewCli(cfg)

	tests := []struct {
		name string
		resp api.SuccessResponse
		txt  string
		call func() (string, error)
	}{
		{
			name: "ShowWhitelist",
			resp: api.SuccessResponse{Result: ips},
			txt:  "--=== Whitelist ===--\n" + strings.Join(ips, "\n") + "\n",
			call: func() (string, error) { return appCli.ShowWhitelist() },
		},
		{
			name: "ShowWhitelist no ips",
			resp: api.SuccessResponse{Result: []string{}},
			txt:  "--=== Whitelist ===--\nEmpty\n",
			call: func() (string, error) { return appCli.ShowWhitelist() },
		},
		{
			name: "ShowBlacklist",
			resp: api.SuccessResponse{Result: ips},
			txt:  "--=== Blacklist ===--\n" + strings.Join(ips, "\n") + "\n",
			call: func() (string, error) { return appCli.ShowBlacklist() },
		},
		{
			name: "ShowBlacklist no ips",
			resp: api.SuccessResponse{Result: []string{}},
			txt:  "--=== Blacklist ===--\nEmpty\n",
			call: func() (string, error) { return appCli.ShowBlacklist() },
		},
		{
			name: "Add blacklist",
			resp: api.SuccessResponse{Result: ips},
			txt:  "OK",
			call: func() (string, error) { return appCli.AppendBlacklist(ip) },
		},
		{
			name: "Remove blacklist",
			resp: api.SuccessResponse{Result: ips},
			txt:  "OK",
			call: func() (string, error) { return appCli.RemoveBlacklist(ip) },
		},
		{
			name: "Add whitelist",
			resp: api.SuccessResponse{Result: ips},
			txt:  "OK",
			call: func() (string, error) { return appCli.AppendWhitelist(ip) },
		},
		{
			name: "Remove whitelist",
			resp: api.SuccessResponse{Result: ips},
			txt:  "OK",
			call: func() (string, error) { return appCli.RemoveWhitelist(ip) },
		},
		{
			name: "Drop bucket login",
			resp: api.SuccessResponse{Result: ips},
			txt:  "OK",
			call: func() (string, error) { return appCli.DropBucketByLogin("Ivan") },
		},
		{
			name: "Drop bucket pwd",
			resp: api.SuccessResponse{Result: ips},
			txt:  "OK",
			call: func() (string, error) { return appCli.DropBucketByPwd("123456") },
		},
		{
			name: "Drop bucket IP",
			resp: api.SuccessResponse{Result: ips},
			txt:  "OK",
			call: func() (string, error) { return appCli.DropBucketByIP(ip) },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := newAPIClientMock(http.StatusOK, test.resp.JSON(), nil)
			appCli.client.setHTTPClient(mock)
			txt, err := test.call()
			require.NoError(t, err)
			require.Equal(t, test.txt, txt)
		})
	}
}

func TestCli_FailedRequest(t *testing.T) {
	ip := "172.168.10.5"
	cfg := config.AppConfig{Addr: "0.0.0.0:8081"}
	appCli := NewCli(cfg)
	expErr := errors.New("my test error")
	tests := []struct {
		name string
		resp api.ErrorResponse
		err  error
		call func() (string, error)
	}{
		{
			name: "ShowWhitelist",
			resp: api.ErrorResponse{ErrMessage: "some error"},
			call: func() (string, error) { return appCli.ShowWhitelist() },
			err:  expErr,
		},
		{
			name: "ShowBlacklist",
			resp: api.ErrorResponse{ErrMessage: "some error"},
			call: func() (string, error) { return appCli.ShowBlacklist() },
			err:  expErr,
		},
		{
			name: "Add blacklist",
			resp: api.ErrorResponse{ErrMessage: "some error"},
			call: func() (string, error) { return appCli.AppendBlacklist(ip) },
			err:  expErr,
		},
		{
			name: "Remove blacklist",
			resp: api.ErrorResponse{ErrMessage: "some error"},
			call: func() (string, error) { return appCli.RemoveBlacklist(ip) },
			err:  expErr,
		},
		{
			name: "Add whitelist",
			resp: api.ErrorResponse{ErrMessage: "some error"},
			call: func() (string, error) { return appCli.AppendWhitelist(ip) },
			err:  expErr,
		},
		{
			name: "Remove whitelist",
			resp: api.ErrorResponse{ErrMessage: "some error"},
			call: func() (string, error) { return appCli.RemoveWhitelist(ip) },
			err:  expErr,
		},
		{
			name: "Drop bucket login",
			resp: api.ErrorResponse{ErrMessage: "some error"},
			call: func() (string, error) { return appCli.DropBucketByLogin("Ivan") },
			err:  expErr,
		},
		{
			name: "Drop bucket pwd",
			resp: api.ErrorResponse{ErrMessage: "some error"},
			call: func() (string, error) { return appCli.DropBucketByPwd("123456") },
			err:  expErr,
		},
		{
			name: "Drop bucket IP",
			resp: api.ErrorResponse{ErrMessage: "some error"},
			call: func() (string, error) { return appCli.DropBucketByIP(ip) },
			err:  expErr,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := newAPIClientMock(http.StatusBadRequest, test.resp.JSON(), test.err)
			appCli.client.setHTTPClient(mock)
			_, err := test.call()
			require.ErrorIs(t, err, test.err)
		})
	}
}
