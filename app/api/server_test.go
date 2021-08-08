package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
)

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
