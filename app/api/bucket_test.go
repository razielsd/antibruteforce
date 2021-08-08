package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAbfAPI_Drop_EmptyKey(t *testing.T) {
	api := createServer()
	tests := []struct {
		name   string
		action func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:   "DropLogin",
			action: api.handlerDropLogin,
		},
		{
			name:   "DropPwd",
			action: api.handlerDropPasswd,
		},
		{
			name:   "DropIP",
			action: api.handlerDropIP,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w, r := createPostReqAndWriter("")
			test.action(w, r)
			require.Equal(t, http.StatusBadRequest, w.Code)

			exp := ErrorResponse{
				Code:       ErrCodeEmptyParam,
				ErrMessage: "empty param: require param key",
			}
			require.JSONEq(t, exp.JSON(), w.Body.String())
		})
	}
}

func TestAbfAPI_DropLogin_ExistsKey(t *testing.T) {
	api := createServer()
	allowParam := createAllowParam("Ivan", "123456", "192.168.1.71")
	isDisallow := fillServerForParam(t, api, allowParam)
	require.True(t, isDisallow)
	w, r := createPostReqAndWriter("key=Ivan")
	api.handlerDropLogin(w, r)
	assertResponseOK(t, w)
	w, r = createPostReqAndWriter(allowParam)
	api.handlerUserAllow(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	resp := &SuccessResponse{
		Result: AllowResult{
			Login: true,
		},
	}
	require.JSONEq(t, resp.JSON(), w.Body.String())
}

func TestAbfAPI_DropPwd_ExistsKey(t *testing.T) {
	api := createServer()
	allowParam := createAllowParam("Ivan", "123456", "192.168.1.71")
	isDisallow := fillServerForParam(t, api, allowParam)
	require.True(t, isDisallow)
	w, r := createPostReqAndWriter("key=123456")
	api.handlerDropPasswd(w, r)
	assertResponseOK(t, w)
	w, r = createPostReqAndWriter(allowParam)
	api.handlerUserAllow(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	resp := &SuccessResponse{
		Result: AllowResult{
			Pwd: true,
		},
	}
	require.JSONEq(t, resp.JSON(), w.Body.String())
}

func TestAbfAPI_DropIP_ExistsKey(t *testing.T) {
	api := createServer()
	allowParam := createAllowParam("Ivan", "123456", "192.168.1.71")
	isDisallow := fillServerForParam(t, api, allowParam)
	require.True(t, isDisallow)
	w, r := createPostReqAndWriter("key=192.168.1.71")
	api.handlerDropIP(w, r)
	assertResponseOK(t, w)
	w, r = createPostReqAndWriter(allowParam)
	api.handlerUserAllow(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	resp := &SuccessResponse{
		Result: AllowResult{
			IP: true,
		},
	}
	require.JSONEq(t, resp.JSON(), w.Body.String())
}

func fillServerForParam(t *testing.T, api *AbfAPI, params string) bool {
	for i := 0; i <= testAPIRate; i++ {
		w, r := createPostReqAndWriter(params)
		api.handlerUserAllow(w, r)
		require.Equal(t, http.StatusOK, w.Code)
		resp := struct {
			Result map[string]bool
		}{
			Result: make(map[string]bool),
		}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		if resp.Result["can-auth"] == false {
			require.False(t, resp.Result["login-check"])
			require.False(t, resp.Result["pwd-check"])
			require.False(t, resp.Result["ip-check"])
			return true
		}
	}
	return false
}
