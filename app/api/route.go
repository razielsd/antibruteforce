package api

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (a *AbfAPI) addRoute(r *mux.Router) {
	r.HandleFunc("/api/user/allow", a.UserAllow).Methods("POST", "GET")

	r.HandleFunc("/api/whitelist", a.GetWhitelist).Methods("GET")
	r.HandleFunc("/api/whitelist/add", a.AppendWhitelist).Methods("POST")
	r.HandleFunc("/api/whitelist/remove", a.RemoveWhitelist).Methods("POST")

	r.HandleFunc("/api/blacklist", a.GetBlacklist).Methods("GET")
	r.HandleFunc("/api/blacklist/add", a.AppendBlacklist).Methods("POST")
	r.HandleFunc("/api/blacklist/remove", a.RemoveBlacklist).Methods("POST")

	r.HandleFunc("/api/bucket/drop/login", a.DropLogin).Methods("POST")
	r.HandleFunc("/api/bucket/drop/pwd", a.DropPasswd).Methods("POST")
	r.HandleFunc("/api/bucket/drop/ip", a.DropIP).Methods("POST")

	r.HandleFunc("/health/liveness", a.ActionHealthProbe).Methods("GET")
	r.HandleFunc("/health/readiness", a.ActionHealthProbe).Methods("GET")
	r.Path("/metrics").Handler(promhttp.Handler())
}
