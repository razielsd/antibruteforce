// +build integration

package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/razielsd/antibruteforce/app/api"
	"github.com/stretchr/testify/require"
)

const testHost = "127.0.0.1:51213"

func TestClientAPI_GetBlacklist(t *testing.T) {
	expectedIP := []string{"192.168.1.71", "10.1.10.0/24"}
	cancel := startTestServer(
		t,
		map[string]api.SuccessResponse{"/api/blacklist": {Result: expectedIP}},
	)
	defer cancel()
	client := newClientAPI(testHost)
	ips, err := client.getBlacklist()
	require.NoError(t, err)
	require.Equal(t, expectedIP, ips)
}

func TestClientAPI_GetWhitelist(t *testing.T) {
	expectedIP := []string{"192.168.1.71", "10.1.10.0/24"}
	cancel := startTestServer(
		t,
		map[string]api.SuccessResponse{"/api/whitelist": {Result: expectedIP}},
	)
	defer cancel()
	client := newClientAPI(testHost)
	ips, err := client.getWhitelist()
	require.NoError(t, err)
	require.Equal(t, expectedIP, ips)
}

func TestModifyBWList(t *testing.T) {
	client := newClientAPI(testHost)

	tests := []struct {
		name    string
		uri     string
		handler func(clientIP string) error
	}{
		{
			name:    "add whitelist",
			uri:     "/api/whitelist/add",
			handler: client.appendWhitelist,
		},
		{
			name:    "remove whitelist",
			uri:     "/api/whitelist/remove",
			handler: client.removeWhitelist,
		},
		{
			name:    "add blacklist",
			uri:     "/api/blacklist/add",
			handler: client.appendBlacklist,
		},
		{
			name:    "remove blacklist",
			uri:     "/api/blacklist/remove",
			handler: client.removeBlacklist,
		},
		{
			name:    "drop bucket by login",
			uri:     "/api/bucket/drop/login",
			handler: client.dropBucketByLogin,
		},
		{
			name:    "drop bucket by password",
			uri:     "/api/bucket/drop/pwd",
			handler: client.dropBucketByPasswd,
		},
		{
			name:    "drop bucket by ip",
			uri:     "/api/bucket/drop/ip",
			handler: client.dropBucketByIP,
		},
	}

	handlers := make(map[string]api.SuccessResponse)

	for _, test := range tests {
		handlers[test.uri] = api.SuccessResponse{
			Result: api.NewSuccessOK(),
		}
	}
	cancel := startTestServer(t, handlers)
	defer cancel()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.handler("192.168.1.99")
			require.NoError(t, err)
		})
	}
}

func startTestServer(t *testing.T, handlers map[string]api.SuccessResponse) context.CancelFunc {
	r := mux.NewRouter()
	for path := range handlers {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(handlers[r.RequestURI])
			require.NoError(t, err)
		}
		r.HandleFunc(path, handler)
	}
	srv := http.Server{
		Addr:    testHost,
		Handler: r,
	}
	ctx, cancel := context.WithCancel(context.Background())
	go srv.ListenAndServe()
	go func() {
		<-ctx.Done()
		st, stCancel := context.WithTimeout(context.Background(), time.Second)
		defer stCancel()
		srv.Shutdown(st)
	}()
	return cancel
}

func stopTestServer(t *testing.T, srv http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	require.NoError(t, err)
}
