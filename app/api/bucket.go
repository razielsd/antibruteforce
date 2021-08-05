package api

import (
	"net/http"

	"go.uber.org/zap"
)

func (a *AbfAPI) DropLogin(w http.ResponseWriter, r *http.Request) {
	form, ok := a.getForm(w, r, []string{"key"})
	if !ok {
		return
	}
	a.loginLimiter.Remove(form["key"])
	a.log.Error("Drop bucket for login", zap.String("login", form["key"]))
	a.sendResult(w, NewSuccessOK())
}

func (a *AbfAPI) DropPasswd(w http.ResponseWriter, r *http.Request) {
	form, ok := a.getForm(w, r, []string{"key"})
	if !ok {
		return
	}
	pwd := encodePwd(form["key"])
	a.pwdLimiter.Remove(pwd)
	a.log.Error("Drop bucket for password")
	a.sendResult(w, NewSuccessOK())
}

func (a *AbfAPI) DropIP(w http.ResponseWriter, r *http.Request) {
	form, ok := a.getForm(w, r, []string{"key"})
	if !ok {
		return
	}
	a.ipLimiter.Remove(form["key"])
	a.log.Error("Drop bucket for IP", zap.String("IP", form["key"]))
	a.sendResult(w, NewSuccessOK())
}
