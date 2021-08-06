package api

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/razielsd/antibruteforce/app/iptable"
)

func (a *AbfAPI) GetWhitelist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.sendResult(w, a.whitelist.GetAll())
}

func (a *AbfAPI) GetBlacklist(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.sendResult(w, a.blacklist.GetAll())
}

func (a *AbfAPI) AppendWhitelist(w http.ResponseWriter, r *http.Request) {
	a.bwlistAdd("Whitelist", a.whitelist, w, r)
}

func (a *AbfAPI) AppendBlacklist(w http.ResponseWriter, r *http.Request) {
	a.bwlistAdd("Blacklist", a.blacklist, w, r)
}

func (a *AbfAPI) bwlistAdd(srcName string, srcTable *iptable.IPTable, w http.ResponseWriter, r *http.Request) {
	form, ok := a.getForm(w, r, []string{"ip"})
	if !ok {
		return
	}
	ip := form["ip"]
	err := srcTable.Add(ip)
	if err != nil {
		a.log.Error(srcName+" - unable add ip", zap.Error(err))
		a.sendError(w, ErrCodeUnableCheckIP, "Unable add ip", err)
		return
	}
	a.log.Info(srcName+" added ip", zap.String("IP", ip))
	a.sendResult(w, NewSuccessOK())
}

func (a *AbfAPI) RemoveWhitelist(w http.ResponseWriter, r *http.Request) {
	a.bwlistRemove("Whitelist", a.whitelist, w, r)
}

func (a *AbfAPI) RemoveBlacklist(w http.ResponseWriter, r *http.Request) {
	a.bwlistRemove("Blacklist", a.blacklist, w, r)
}

func (a *AbfAPI) bwlistRemove(srcName string, srcTable *iptable.IPTable, w http.ResponseWriter, r *http.Request) {
	form, ok := a.getForm(w, r, []string{"ip"})
	if !ok {
		return
	}
	ip := form["ip"]
	err := srcTable.Remove(ip)
	if err != nil {
		a.log.Error(srcName+" - unable remove ip", zap.Error(err))
		a.sendError(w, ErrCodeUnableCheckIP, "Unable remove ip", err)
		return
	}
	a.log.Info(srcName+" removed ip", zap.String("IP", ip))
	a.sendResult(w, NewSuccessOK())
}
