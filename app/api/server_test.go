package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/razielsd/antibruteforce/app/config"
	"github.com/razielsd/antibruteforce/app/iptable"
	"github.com/razielsd/antibruteforce/app/logger"
)

func TestLoadIPList(t *testing.T) {
	abf := createServer()
	table := iptable.NewIPTable()
	err := abf.loadIPList(table, []string{"192.168.1.1", "", "10.10.1.0/24", "10.1000"})
	expIP := []string{"192.168.1.1", "10.10.1.0/24"}
	sort.Strings(expIP)
	require.Equal(t, expIP, table.GetAll())
	require.Len(t, err, 1)
}

func createServer() *AbfAPI {
	cfg := config.AppConfig{}
	cfg.RateLogin = testAPIRate
	cfg.RatePwd = testAPIRate
	cfg.RateIP = testAPIRate

	l, _ := logger.GetLogger(cfg)
	api, _ := NewAbfAPI(cfg, l)
	return api
}

func createAllowParam(login, pwd, ip string) string {
	data := url.Values{}
	data.Set("login", login)
	data.Set("pwd", pwd)
	data.Set("ip", ip)
	return data.Encode()
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
