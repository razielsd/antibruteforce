package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAbfAPI_ActionHealthProbe(t *testing.T) {
	w, r := createGetReqAndWriter()
	api := createServer()
	api.handlerHealthProbe(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	exp := SuccessResponse{Result: NewSuccessOK()}
	require.JSONEq(t, exp.JSON(), w.Body.String())
}
