package cli

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/razielsd/antibruteforce/app/api"
)

const testHost = "127.0.0.1:51215"

func TestApiClient_GetWhitelist(t *testing.T) {
	expIps := []string{"192.168.1.71", "10.10.1.15"}
	resp := api.SuccessResponse{Result: expIps}
	mock := newAPIClientMock(http.StatusOK, resp.JSON(), nil)
	client := newClientAPI(testHost)
	client.setHTTPClient(mock)
	ips, err := client.getWhitelist()
	require.NoError(t, err)
	require.Equal(t, expIps, ips)
}

func TestApiClient_GetBlacklist(t *testing.T) {
	expIps := []string{"192.168.1.72", "10.10.1.151"}
	resp := api.SuccessResponse{Result: expIps}
	mock := newAPIClientMock(http.StatusOK, resp.JSON(), nil)
	client := newClientAPI(testHost)
	client.setHTTPClient(mock)
	ips, err := client.getBlacklist()
	require.NoError(t, err)
	require.Equal(t, expIps, ips)
}

func TestApiClient_UpdateRemove_Success(t *testing.T) {
	ip := "172.168.10.5"
	tests := []struct {
		name string
		resp api.SuccessResponse
		exp  interface{}
		call func(c *clientAPI) error
	}{
		{
			name: "Add blacklist",
			resp: api.SuccessResponse{Result: api.NewSuccessOK()},
			call: func(c *clientAPI) error { return c.appendBlacklist(ip) },
		},
		{
			name: "Remove blacklist",
			resp: api.SuccessResponse{Result: api.NewSuccessOK()},
			call: func(c *clientAPI) error { return c.removeBlacklist(ip) },
		},
		{
			name: "Add whitelist",
			resp: api.SuccessResponse{Result: api.NewSuccessOK()},
			call: func(c *clientAPI) error { return c.appendWhitelist(ip) },
		},
		{
			name: "Remove whitelist",
			resp: api.SuccessResponse{Result: api.NewSuccessOK()},
			call: func(c *clientAPI) error { return c.removeWhitelist(ip) },
		},
		{
			name: "Drop bucket login",
			resp: api.SuccessResponse{Result: api.NewSuccessOK()},
			call: func(c *clientAPI) error { return c.dropBucketByLogin("Ivan") },
		},
		{
			name: "Drop bucket pwd",
			resp: api.SuccessResponse{Result: api.NewSuccessOK()},
			call: func(c *clientAPI) error { return c.dropBucketByPasswd("123456") },
		},
		{
			name: "Drop bucket IP",
			resp: api.SuccessResponse{Result: api.NewSuccessOK()},
			call: func(c *clientAPI) error { return c.dropBucketByIP(ip) },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := newAPIClientMock(http.StatusOK, test.resp.JSON(), nil)
			client := newClientAPI(testHost)
			client.setHTTPClient(mock)
			err := test.call(client)
			require.NoError(t, err)
		})
	}
}

func TestSendUpdate_ErrorResponse(t *testing.T) {
	resp := api.ErrorResponse{
		ErrMessage: "was error here",
		Code:       1,
	}
	mock := newAPIClientMock(http.StatusBadRequest, resp.JSON(), nil)
	client := newClientAPI(testHost)
	client.setHTTPClient(mock)
	err := client.sendUpdate("/api/blacklist/add", "ip", "10.20.30.40")
	require.Error(t, err)
	require.Equal(t, resp.ErrMessage, err.Error())
}

func TestSendUpdate_ResponseWithBadData(t *testing.T) {
	resp := "{a:}"
	mock := newAPIClientMock(http.StatusBadRequest, resp, nil)
	client := newClientAPI(testHost)
	client.setHTTPClient(mock)
	err := client.sendUpdate("/api/blacklist/add", "ip", "10.20.30.40")
	require.Error(t, err)
}
