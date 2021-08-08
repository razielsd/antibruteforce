package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
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
		loginLimiter: reqlimiter.NewReqLimiter(reqlimiter.NewLimiterConfig(cfg.RateLogin)),
		pwdLimiter:   reqlimiter.NewReqLimiter(reqlimiter.NewLimiterConfig(cfg.RatePwd)),
		ipLimiter:    reqlimiter.NewReqLimiter(reqlimiter.NewLimiterConfig(cfg.RateIP)),
		whitelist:    iptable.NewIPTable(),
		blacklist:    iptable.NewIPTable(),
	}
	return api, nil
}

func (a *AbfAPI) Run(ctx context.Context) {
	a.initBWList()
	r := mux.NewRouter()
	a.addRoute(r)
	srv := a.createServer(r)
	a.startServer(srv)
	<-ctx.Done()
	a.stopServer(srv)
	a.log.Info("server stopped")
}

func (a *AbfAPI) startServer(srv *http.Server) {
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
}

func (a *AbfAPI) createServer(r *mux.Router) *http.Server {
	return &http.Server{
		Addr:         a.cfg.Addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      newHTTPLog(r, a.log),
	}
}

func (a *AbfAPI) stopServer(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		a.log.Error("error during shutdown", zap.Error(err))
	}
}

func (a *AbfAPI) initBWList() {
	errLoad := a.loadIPList(a.whitelist, a.cfg.Whitelist)
	a.logLoadedIPList(errLoad, "cannot add ip to whitelist from env")
	errLoad = a.loadIPList(a.blacklist, a.cfg.Blacklist)
	a.logLoadedIPList(errLoad, "cannot add ip to whitelist from env")
}

func (a *AbfAPI) loadIPList(table *iptable.IPTable, ips []string) map[string]error {
	result := make(map[string]error)
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			continue
		}
		err := table.Add(ip)
		if err != nil {
			result[ip] = err
		}
	}
	return result
}

func (a *AbfAPI) logLoadedIPList(errLoad map[string]error, message string) {
	for ip, err := range errLoad {
		a.log.Error(
			message,
			zap.Error(err),
			zap.String("IP", ip),
		)
	}
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
