package cli

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApiClientMock_Do_Success(t *testing.T) {
	expCode := http.StatusOK
	expBody := "hello world success"
	mock := newAPIClientMock(expCode, expBody, nil)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "", nil)
	require.NoError(t, err)

	resp, err := mock.Do(req)
	require.NoError(t, err)
	require.Equal(t, expCode, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, expBody, string(body))
}

func TestApiClientMock_Do_Error(t *testing.T) {
	expCode := http.StatusBadRequest
	expBody := "hello world with error"
	expErr := errors.New("hello error")
	mock := newAPIClientMock(expCode, expBody, expErr)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "", nil)
	require.NoError(t, err)

	resp, err := mock.Do(req)
	require.ErrorIs(t, expErr, err)
	require.Equal(t, expCode, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, expBody, string(body))
}

func TestNewApiClientMock(t *testing.T) {
	expCode := 200
	expBody := "hello world"
	mock := newAPIClientMock(expCode, expBody, nil)
	require.Equal(t, expCode, mock.statusCode)
	require.Equal(t, expBody, mock.body)
}
