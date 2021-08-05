package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/razielsd/antibruteforce/app/config"
	"github.com/razielsd/antibruteforce/app/iptable"
	"github.com/razielsd/antibruteforce/app/reqlimiter"
)

const (
	ErrCodeUnableParseRequest = 1
	ErrCodeUnableCheckIP      = 2
	ErrCodeEmptyParam         = 3
)

type AbfAPI struct {
	cfg          config.AppConfig
	log          *zap.Logger
	loginLimiter *reqlimiter.ReqLimiter
	pwdLimiter   *reqlimiter.ReqLimiter
	ipLimiter    *reqlimiter.ReqLimiter
	whitelist    *iptable.IPTable
	blacklist    *iptable.IPTable
}

func NewAbfAPI(cfg config.AppConfig, logger *zap.Logger) (*AbfAPI, error) {
	api := &AbfAPI{
		cfg:          cfg,
		log:          logger,
		loginLimiter: reqlimiter.NewReqLimiter(cfg.RateLogin),
		pwdLimiter:   reqlimiter.NewReqLimiter(cfg.RatePwd),
		ipLimiter:    reqlimiter.NewReqLimiter(cfg.RateIP),
		whitelist:    iptable.NewIPTable(),
		blacklist:    iptable.NewIPTable(),
	}
	return api, nil
}

func (a *AbfAPI) Run() {
	r := mux.NewRouter()

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

	r.Path("/metrics").Handler(promhttp.Handler())

	srv := &http.Server{
		Addr:         a.cfg.Addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      newHTTPLog(r, a.log),
	}

	go func() {
		a.log.Info("start server", zap.String("host", a.cfg.Addr))
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				a.log.Info("server closed")
			} else {
				a.log.Error("error in server", zap.Error(err))
			}
		}
	}()

	c := make(chan os.Signal, 1)
	// @todo: signal.NotifyContext()
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		a.log.Error("error on shutdown", zap.Error(err))
	}
	a.log.Info("shutting down")
}

func (a *AbfAPI) getForm(w http.ResponseWriter, r *http.Request, params []string) (map[string]string, bool) {
	if err := r.ParseForm(); err != nil {
		a.sendError(w, ErrCodeUnableParseRequest, "cannot parse post params", err)
		return nil, false
	}
	result := make(map[string]string)
	for _, key := range params {
		v := r.Form.Get(key)
		if v == "" {
			a.sendError(w, ErrCodeEmptyParam, "empty param", errors.New("require param "+key))
			return nil, false
		}
		result[key] = v
	}
	return result, true
}

func (a *AbfAPI) sendResult(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(SuccessResponse{Result: data})
	if err != nil {
		a.log.Error("Unable encode response", zap.Error(err))
	}
}

func (a *AbfAPI) sendError(w http.ResponseWriter, code int, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	encErr := json.NewEncoder(w).Encode(
		ErrorResponse{
			ErrMessage: fmt.Sprintf("%s: %s", message, err),
			Code:       code,
		})
	if encErr != nil {
		a.log.Error("Error encode response", zap.Error(err))
	}
}
