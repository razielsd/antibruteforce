package api

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (a *AbfAPI) addRoute(r *mux.Router) {
	r.HandleFunc("/api/user/allow", a.handlerUserAllow).Methods("POST", "GET")

	r.HandleFunc("/api/whitelist", a.handlerGetWhitelist).Methods("GET")
	r.HandleFunc("/api/whitelist/add", a.handlerAppendWhitelist).Methods("POST")
	r.HandleFunc("/api/whitelist/remove", a.handlerRemoveWhitelist).Methods("POST")

	r.HandleFunc("/api/blacklist", a.handlerGetBlacklist).Methods("GET")
	r.HandleFunc("/api/blacklist/add", a.handlerAppendBlacklist).Methods("POST")
	r.HandleFunc("/api/blacklist/remove", a.handlerRemoveBlacklist).Methods("POST")

	r.HandleFunc("/api/bucket/drop/login", a.handlerDropLogin).Methods("POST")
	r.HandleFunc("/api/bucket/drop/pwd", a.handlerDropPasswd).Methods("POST")
	r.HandleFunc("/api/bucket/drop/ip", a.handlerDropIP).Methods("POST")

	r.HandleFunc("/health/liveness", a.handlerHealthProbe).Methods("GET")
	r.HandleFunc("/health/readiness", a.handlerHealthProbe).Methods("GET")
	r.Path("/metrics").Handler(promhttp.Handler())
}
