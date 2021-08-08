package api

import (
	"net/http"

	"go.uber.org/zap"
)

type AllowResult struct {
	CanAuth   bool `json:"can-auth"`
	Login     bool `json:"login-check"`
	Pwd       bool `json:"pwd-check"`
	IP        bool `json:"ip-check"`
	Whitelist bool `json:"whitelist-check"`
	Blacklist bool `json:"blacklist-check"`
}

func (a *AbfAPI) handlerUserAllow(w http.ResponseWriter, r *http.Request) {
	form, ok := a.getForm(w, r, []string{"login", "pwd", "ip"})
	if !ok {
		return
	}
	login := form["login"]
	pwd := form["pwd"]
	ip := form["ip"]

	pwd = encodePwd(pwd)
	res := AllowResult{}

	isAllow, err := a.whitelist.Contains(ip)
	if err != nil {
		a.sendError(w, ErrCodeUnableCheckIP, "ip is "+ip, err)
		return
	}
	if isAllow {
		res.Whitelist = true
		a.userAllowWriteResult(w, res, login, ip)
		return
	}

	disallow, err := a.blacklist.Contains(ip)
	if err != nil {
		a.sendError(w, ErrCodeUnableCheckIP, "ip", err)
		return
	}
	if disallow {
		res.Blacklist = true
		a.userAllowWriteResult(w, res, login, ip)
		return
	}

	res.Login = a.loginLimiter.Allow(login)
	res.Pwd = a.pwdLimiter.Allow(pwd)
	res.IP = a.ipLimiter.Allow(ip)
	a.userAllowWriteResult(w, res, login, ip)
}

func (a *AbfAPI) userAllowWriteResult(w http.ResponseWriter, res AllowResult, login, ip string) {
	res.Build()
	if !res.CanAuth {
		a.log.Info(
			"Block allow",
			zap.String("login", login),
			zap.String("ip", ip),
			zap.Bool("login-check", res.Login),
			zap.Bool("pwd-check", res.Pwd),
			zap.Bool("ip-check", res.IP),
			zap.Bool("whitelist-check", res.Whitelist),
			zap.Bool("blacklist-check", res.Blacklist),
		)
	}
	a.sendResult(w, res)
}

func (a *AllowResult) Build() {
	a.CanAuth = false
	switch {
	case a.Whitelist:
		a.CanAuth = true
	case a.Blacklist:
		a.CanAuth = false
	case a.Login && a.Pwd && a.IP:
		a.CanAuth = true
	}
}
