package api

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

const clientIP = "192.168.1.71"

func TestAbfAPI_GetBlacklist(t *testing.T) {
	api := createServer()
	err := api.blacklist.Add(clientIP)
	require.NoError(t, err)

	w, r := createGetReqAndWriter()
	api.GetBlacklist(w, r)
	assertResponseContainsIP(t, w, []string{clientIP})
}

func TestAbfAPI_GetWhitelist(t *testing.T) {
	api := createServer()
	err := api.whitelist.Add(clientIP)
	require.NoError(t, err)
	w, r := createGetReqAndWriter()
	api.GetWhitelist(w, r)

	assertResponseContainsIP(t, w, []string{clientIP})
}

func TestAbfAPI_AppendWhitelist_ValidIP_Success(t *testing.T) {
	api := createServer()
	w, r := createPostReqAndWriter("ip=" + clientIP)
	api.AppendWhitelist(w, r)
	assertResponseOK(t, w)

	w, r = createGetReqAndWriter()
	api.GetWhitelist(w, r)
	assertResponseContainsIP(t, w, []string{clientIP})
}

func TestAbfAPI_AppendWhitelist_EmptyIP_Error(t *testing.T) {
	api := createServer()
	w, r := createPostReqAndWriter("")
	api.AppendWhitelist(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)

	exp := ErrorResponse{
		Code:       ErrCodeEmptyParam,
		ErrMessage: "empty param: require param ip",
	}
	require.JSONEq(t, exp.JSON(), w.Body.String())
}

func TestAbfAPI_AppendBlacklist_ValidIP_Success(t *testing.T) {
	api := createServer()
	w, r := createPostReqAndWriter("ip=" + clientIP)
	api.AppendBlacklist(w, r)
	assertResponseOK(t, w)

	w, r = createGetReqAndWriter()
	api.GetBlacklist(w, r)
	assertResponseContainsIP(t, w, []string{clientIP})
}

func TestAbfAPI_AppendBlacklist_EmptyIP_Error(t *testing.T) {
	api := createServer()
	w, r := createPostReqAndWriter("")
	api.AppendBlacklist(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)

	exp := ErrorResponse{
		Code:       ErrCodeEmptyParam,
		ErrMessage: "empty param: require param ip",
	}
	require.JSONEq(t, exp.JSON(), w.Body.String())
}

func TestAbfAPI_RemoveWhitelist_ValidIP_Success(t *testing.T) {
	api := createServer()

	w, r := createGetReqAndWriter()
	api.GetWhitelist(w, r)
	assertResponseContainsIP(t, w, []string{})

	w, r = createPostReqAndWriter("ip=" + clientIP)
	api.AppendWhitelist(w, r)
	assertResponseOK(t, w)

	w, r = createGetReqAndWriter()
	api.GetWhitelist(w, r)
	assertResponseContainsIP(t, w, []string{clientIP})

	w, r = createPostReqAndWriter("ip=" + clientIP)
	api.RemoveWhitelist(w, r)
	assertResponseOK(t, w)

	w, r = createGetReqAndWriter()
	api.GetWhitelist(w, r)
	assertResponseContainsIP(t, w, []string{})
}

func TestAbfAPI_RemoveWhitelist_EmptyIP_Error(t *testing.T) {
	api := createServer()
	w, r := createPostReqAndWriter("")
	api.RemoveWhitelist(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)

	exp := ErrorResponse{
		Code:       ErrCodeEmptyParam,
		ErrMessage: "empty param: require param ip",
	}
	require.JSONEq(t, exp.JSON(), w.Body.String())
}

func TestAbfAPI_RemoveBlacklist_ValidIP_Success(t *testing.T) {
	api := createServer()

	w, r := createGetReqAndWriter()
	api.GetBlacklist(w, r)
	assertResponseContainsIP(t, w, []string{})

	w, r = createPostReqAndWriter("ip=" + clientIP)
	api.AppendBlacklist(w, r)
	assertResponseOK(t, w)

	w, r = createGetReqAndWriter()
	api.GetBlacklist(w, r)
	assertResponseContainsIP(t, w, []string{clientIP})

	w, r = createPostReqAndWriter("ip=" + clientIP)
	api.RemoveBlacklist(w, r)
	assertResponseOK(t, w)

	w, r = createGetReqAndWriter()
	api.GetBlacklist(w, r)
	assertResponseContainsIP(t, w, []string{})
}

func TestAbfAPI_RemoveBlacklist_EmptyIP_Error(t *testing.T) {
	api := createServer()
	w, r := createPostReqAndWriter("")
	api.RemoveBlacklist(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)

	exp := ErrorResponse{
		Code:       ErrCodeEmptyParam,
		ErrMessage: "empty param: require param ip",
	}
	require.JSONEq(t, exp.JSON(), w.Body.String())
}

func assertResponseOK(t *testing.T, w *httptest.ResponseRecorder) {
	require.Equal(t, http.StatusOK, w.Code)
	exp := SuccessResponse{
		Result: NewSuccessOK(),
	}
	require.JSONEq(t, exp.JSON(), w.Body.String())
}

func assertResponseContainsIP(t *testing.T, w *httptest.ResponseRecorder, ips []string) {
	require.Equal(t, http.StatusOK, w.Code)
	sort.Strings(ips)
	exp := SuccessResponse{
		Result: ips,
	}
	require.JSONEq(t, exp.JSON(), w.Body.String())
}
