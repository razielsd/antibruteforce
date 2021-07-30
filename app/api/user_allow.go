package api

import (
	"encoding/json"
	"net/http"
)

type allowResult struct {
	CanAuth bool `json:"can-auth"`
	Login   bool `json:"login-check"`
	Pwd     bool `json:"pwd-check"`
	IP      bool `json:"ip-check"`
}

func (a *AbfAPI) UserAllow(w http.ResponseWriter, r *http.Request) {
	// login & pwd - will be hashed, now is plain
	login := "ivan"
	pwd := "123123"
	ip := "192.168.1.15"

	res := allowResult{}

	res.Login = a.loginLimiter.Allow(login)
	res.Pwd = a.pwdLimiter.Allow(pwd)
	res.IP = a.ipLimiter.Allow(ip)
	res.Build()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse{Result: res})
}

func (a *allowResult) Build() {
	a.CanAuth = false
	if a.Login && a.Pwd && a.IP {
		a.CanAuth = true
	}
}
