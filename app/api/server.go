package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/razielsd/antibruteforce/app/config"
	"github.com/razielsd/antibruteforce/app/reqlimiter"
	"go.uber.org/zap"
)

type AbfAPI struct {
	cfg          config.AppConfig
	log          *zap.Logger
	loginLimiter *reqlimiter.ReqLimiter
	pwdLimiter   *reqlimiter.ReqLimiter
	ipLimiter    *reqlimiter.ReqLimiter
}

type successResponse struct {
	Result interface{} `json:"result"`
}

func NewAbfAPI(cfg config.AppConfig, logger *zap.Logger) AbfAPI {
	return AbfAPI{
		cfg:          cfg,
		log:          logger,
		loginLimiter: reqlimiter.NewReqLimiter(cfg.RateLogin),
		pwdLimiter:   reqlimiter.NewReqLimiter(cfg.RatePwd),
		ipLimiter:    reqlimiter.NewReqLimiter(cfg.RateIP),
	}
}

func (a *AbfAPI) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/api/user/allow", a.UserAllow)

	srv := &http.Server{
		Addr:         a.cfg.Addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
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
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	a.log.Info("shutting down")
}
