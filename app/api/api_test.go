// +build integration

package api

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var apiHost = "127.0.0.1:51212"

func TestAbfAPI_Run_ReadinessProbe(t *testing.T) {
	cancel := runServer(apiHost)
	defer cancel()
	makeProbe(t, apiHost, "/health/readiness")
}

func TestAbfAPI_Run_LivenessProbe(t *testing.T) {
	cancel := runServer(apiHost)
	defer cancel()
	makeProbe(t, apiHost, "/health/liveness")
}

// Here only check endpoints, really tests in actions tests
func TestAbfAPI_Run_All(t *testing.T) {
	tests := []struct {
		uri     string
		method  string
		body    io.Reader
		expResp SuccessResponse
	}{
		{
			uri:    "/api/user/allow",
			method: http.MethodPost,
			body:   strings.NewReader(createAllowParam("Ivan", "123456", "192.168.1.90")),
			expResp: SuccessResponse{
				Result: AllowResult{
					CanAuth: true,
					Login:   true,
					Pwd:     true,
					IP:      true,
				},
			},
		},
		{
			uri:    "/api/blacklist",
			method: http.MethodGet,
			body:   nil,
			expResp: SuccessResponse{
				Result: []string{},
			},
		},
		{
			uri:    "/api/blacklist/add",
			method: http.MethodPost,
			body:   strings.NewReader("ip=192.168.1.90"),
			expResp: SuccessResponse{
				Result: NewSuccessOK(),
			},
		},
		{
			uri:    "/api/blacklist/remove",
			method: http.MethodPost,
			body:   strings.NewReader("ip=192.168.1.90"),
			expResp: SuccessResponse{
				Result: NewSuccessOK(),
			},
		},
		{
			uri:    "/api/whitelist",
			method: http.MethodGet,
			body:   nil,
			expResp: SuccessResponse{
				Result: []string{},
			},
		},
		{
			uri:    "/api/whitelist/add",
			method: http.MethodPost,
			body:   strings.NewReader("ip=192.168.1.90"),
			expResp: SuccessResponse{
				Result: NewSuccessOK(),
			},
		},
		{
			uri:    "/api/whitelist/remove",
			method: http.MethodPost,
			body:   strings.NewReader("ip=192.168.1.90"),
			expResp: SuccessResponse{
				Result: NewSuccessOK(),
			},
		},
		{
			uri:    "/api/bucket/drop/login",
			method: http.MethodPost,
			body:   strings.NewReader("key=Ivan"),
			expResp: SuccessResponse{
				Result: NewSuccessOK(),
			},
		},
		{
			uri:    "/api/bucket/drop/pwd",
			method: http.MethodPost,
			body:   strings.NewReader("key=123456"),
			expResp: SuccessResponse{
				Result: NewSuccessOK(),
			},
		},
		{
			uri:    "/api/bucket/drop/ip",
			method: http.MethodPost,
			body:   strings.NewReader("key=192.168.1.99"),
			expResp: SuccessResponse{
				Result: NewSuccessOK(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.uri, func(t *testing.T) {
			cancel := runServer(apiHost)
			defer cancel()
			waitReadinessProbe(t, apiHost)
			reqURL := fmt.Sprintf("http://%s%s", apiHost, test.uri)
			req, err := http.NewRequest(
				test.method, reqURL, test.body,
			)
			require.NoError(t, err)
			if test.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
			client := http.Client{}
			resp, err := client.Do(req)
			defer resp.Body.Close()
			require.NoError(t, err)

			body, err := ioutil.ReadAll(resp.Body)
			require.JSONEq(t, test.expResp.JSON(), string(body))
		})
	}
}

func runServer(host string) context.CancelFunc {
	srv := createServer()
	srv.cfg.Addr = host
	ctx, cancel := context.WithCancel(context.Background())
	go srv.Run(ctx)
	return cancel
}

func waitReadinessProbe(t *testing.T, host string) {
	makeProbe(t, host, "/health/readiness")
}

func makeProbe(t *testing.T, host string, uri string) {
	cond := func() bool {
		reqURL := fmt.Sprintf("http://%s%s", host, uri)
		req, _ := http.NewRequestWithContext(
			context.Background(), http.MethodGet, reqURL, nil,
		)
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}
	require.Eventually(t, cond, 20*time.Second, 2*time.Second)
}
