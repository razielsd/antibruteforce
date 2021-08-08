package api

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/razielsd/antibruteforce/app/iptable"
)

func (a *AbfAPI) handlerGetWhitelist(w http.ResponseWriter, r *http.Request) {
	a.sendResult(w, a.whitelist.GetAll())
}

func (a *AbfAPI) handlerGetBlacklist(w http.ResponseWriter, r *http.Request) {
	a.sendResult(w, a.blacklist.GetAll())
}

func (a *AbfAPI) handlerAppendWhitelist(w http.ResponseWriter, r *http.Request) {
	a.bwlistAdd("Whitelist", a.whitelist, w, r)
}

func (a *AbfAPI) handlerAppendBlacklist(w http.ResponseWriter, r *http.Request) {
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

func (a *AbfAPI) handlerRemoveWhitelist(w http.ResponseWriter, r *http.Request) {
	a.bwlistRemove("Whitelist", a.whitelist, w, r)
}

func (a *AbfAPI) handlerRemoveBlacklist(w http.ResponseWriter, r *http.Request) {
	a.bwlistRemove("Blacklist", a.blacklist, w, r)
}

func (a *AbfAPI) bwlistRemove(srcName string, srcTable *iptable.IPTable, w http.ResponseWriter, r *http.Request) {
	form, ok := a.getForm(w, r, []string{"ip"})
	if !ok {
		return
	}
	ip := form["ip"]
	srcTable.Remove(ip)
	a.log.Info(srcName+" removed ip", zap.String("IP", ip))
	a.sendResult(w, NewSuccessOK())
}
