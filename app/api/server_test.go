package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAbfAPI_Run(t *testing.T) {
	srv := createServer()
	host := "127.0.0.1:51212"
	srv.cfg.Addr = host
	ctx, cancel := context.WithCancel(context.Background())
	go srv.Run(ctx)

	cond := func() bool {
		reqURL := fmt.Sprintf("http://%s/health/readiness", host)
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

	require.Eventually(t, cond, 10*time.Second, time.Second)
	defer cancel()
}

func createPostReqAndWriter(paramStr string) (*httptest.ResponseRecorder, *http.Request) {
	r, _ := http.NewRequestWithContext(
		context.Background(), http.MethodPost, "/", strings.NewReader(paramStr),
	)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	return w, r
}

func createGetReqAndWriter() (*httptest.ResponseRecorder, *http.Request) {
	r, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet, "/", nil,
	)
	w := httptest.NewRecorder()
	return w, r
}
