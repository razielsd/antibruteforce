package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/razielsd/antibruteforce/app/config"
	"github.com/razielsd/antibruteforce/app/logger"
)

const testAPIRate = 2

func TestAbfAPI_UserAllow_EmptyBWList_AuthSuccess(t *testing.T) {
	w, r := createPostReqAndWriter(createAllowParam("Ivan", "123456", "192.168.1.71"))
	api := createServer()
	api.handlerUserAllow(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	exp := SuccessResponse{
		Result: AllowResult{
			CanAuth: true,
			Login:   true,
			Pwd:     true,
			IP:      true,
		},
	}
	require.JSONEq(t, exp.JSON(), w.Body.String())
}

func TestAbfAPI_UserAllow_IpInBlacklist_AuthFailed(t *testing.T) {
	w, r := createPostReqAndWriter(createAllowParam("Ivan", "123456", "192.168.1.71"))
	api := createServer()
	err := api.blacklist.Add("192.168.1.71")
	require.NoError(t, err)
	api.handlerUserAllow(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	exp := SuccessResponse{
		Result: AllowResult{
			Blacklist: true,
		},
	}
	require.JSONEq(t, exp.JSON(), w.Body.String())
}

func TestAbfAPI_UserAllow_IpInWhitelist_AuthSuccess(t *testing.T) {
	ip := "192.168.1.71"
	tryTotal := testAPIRate + 2
	api := createServer()
	err := api.whitelist.Add("192.168.1.71")
	require.NoError(t, err)
	exp := SuccessResponse{
		Result: AllowResult{
			CanAuth:   true,
			Whitelist: true,
		},
	}

	for i := 0; i < tryTotal; i++ {
		w, r := createPostReqAndWriter(createAllowParam("Ivan", "123456", ip))
		api.handlerUserAllow(w, r)
		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, exp.JSON(), w.Body.String())
	}
}

func TestAbfAPI_UserAllow_LimitExceed_AuthFailed(t *testing.T) {
	testLogin := "Ivan"
	testPwd := "123456"
	testIP := "192.168.1.17"

	tests := []struct {
		name   string
		login  string
		passwd string
		ip     string
		exp    AllowResult
	}{
		{
			name:   "auth by login exceeded",
			login:  testLogin,
			passwd: testPwd + "!!",
			ip:     testIP + "1",
			exp:    AllowResult{Pwd: true, IP: true},
		},
		{
			name:   "auth by pwd exceeded",
			login:  testLogin + "!",
			passwd: testPwd,
			ip:     testIP + "1",
			exp:    AllowResult{Login: true, IP: true},
		},
		{
			name:   "auth by ip exceeded",
			login:  testLogin + "!",
			passwd: testPwd + "!!",
			ip:     testIP,
			exp:    AllowResult{Login: true, Pwd: true},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := createServer()
			exp := SuccessResponse{
				Result: AllowResult{
					CanAuth: true,
					Login:   true,
					Pwd:     true,
					IP:      true,
				},
			}
			for i := 0; i < testAPIRate; i++ {
				w, r := createPostReqAndWriter(createAllowParam(testLogin, testPwd, testIP))
				api.handlerUserAllow(w, r)
				require.Equal(t, http.StatusOK, w.Code)
				require.JSONEq(t, exp.JSON(), w.Body.String())
			}
			w, r := createPostReqAndWriter(createAllowParam(test.login, test.passwd, test.ip))
			api.handlerUserAllow(w, r)
			require.Equal(t, http.StatusOK, w.Code)
			ej, err := json.Marshal(SuccessResponse{Result: test.exp})
			require.NoError(t, err, "Unable marshal expected json")
			require.JSONEq(t, string(ej), w.Body.String())
		})
	}
}

func TestAllowResult_Build_CanAuthSuccess(t *testing.T) {
	tests := []struct {
		name string
		res  AllowResult
	}{
		{
			name: "whitelist=true",
			res:  AllowResult{Whitelist: true},
		},
		{
			name: "login/pwd/ip=true",
			res:  AllowResult{Login: true, Pwd: true, IP: true},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.res.Build()
			require.True(t, test.res.CanAuth)
		})
	}
}

func TestAbfAPI_UserAllow_WrongIP_Error(t *testing.T) {
	w, r := createPostReqAndWriter(createAllowParam("Ivan", "123456", "192.168.1.711"))
	api := createServer()
	api.handlerUserAllow(w, r)

	require.Equal(t, http.StatusBadRequest, w.Code)

	x := w.Body.String()
	fmt.Println(x)

	exp := ErrorResponse{
		ErrMessage: "ip is 192.168.1.711: invalid ip address:192.168.1.711 is not valid ipv4 address",
		Code:       ErrCodeUnableCheckIP,
	}

	require.JSONEq(t, exp.JSON(), w.Body.String())
}

func TestAbfAPI_UserAllow_EmptyParam_Error(t *testing.T) {
	tests := []struct {
		name  string
		login string
		pwd   string
		ip    string
	}{
		{name: "login", login: "", pwd: "123456", ip: "192.168.1.67"},
		{name: "pwd", login: "Ivan", pwd: "", ip: "192.168.1.67"},
		{name: "ip", login: "Petr", pwd: "123456", ip: ""},
	}
	for _, test := range tests {
		t.Run("empty "+test.name, func(t *testing.T) {
			w, r := createPostReqAndWriter(
				createAllowParam(test.login, test.pwd, test.ip),
			)
			api := createServer()
			api.handlerUserAllow(w, r)
			require.Equal(t, http.StatusBadRequest, w.Code)
			exp := ErrorResponse{
				Code:       ErrCodeEmptyParam,
				ErrMessage: "empty param: require param " + test.name,
			}
			require.JSONEq(t, exp.JSON(), w.Body.String())
		})
	}
}

func TestAllowResult_Build_CanAuthFail(t *testing.T) {
	tests := []struct {
		name string
		res  AllowResult
	}{
		{
			name: "all false",
			res:  AllowResult{},
		},
		{
			name: "blacklist=true",
			res:  AllowResult{Blacklist: true, Login: true, Pwd: true, IP: true},
		},
		{
			name: "single check failed - login",
			res:  AllowResult{Login: false, Pwd: true, IP: true},
		},
		{
			name: "single check failed - pwd",
			res:  AllowResult{Login: true, Pwd: false, IP: true},
		},
		{
			name: "single check failed - ip",
			res:  AllowResult{Login: true, Pwd: true, IP: false},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.res.Build()
			require.False(t, test.res.CanAuth)
		})
	}
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
