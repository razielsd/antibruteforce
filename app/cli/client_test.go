package cli

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/razielsd/antibruteforce/app/api"
)

func TestClientExtractError(t *testing.T) {
	client := newClientAPI(testHost)
	resp := httptest.NewRecorder()
	exp := api.ErrorResponse{
		ErrMessage: "hello error",
		Code:       2,
	}
	_, err := resp.WriteString(exp.JSON())
	require.NoError(t, err)
	apiErr, err := client.extractError(resp.Result())
	require.NoError(t, err)
	require.JSONEq(t, exp.JSON(), apiErr.JSON())
}
