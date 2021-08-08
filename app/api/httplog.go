package api

import (
	"bufio"
	"net"
	"net/http"

	"github.com/felixge/httpsnoop"
	"go.uber.org/zap"
)

type httpLog struct {
	log     *zap.Logger
	handler http.Handler
}

type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func newHTTPLog(handler http.Handler, log *zap.Logger) *httpLog {
	return &httpLog{
		log:     log,
		handler: handler,
	}
}

func (h *httpLog) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	respLogger, wl := newResponseLogger(w)
	h.handler.ServeHTTP(wl, req)
	h.log.Info(
		"Request",
		zap.String("Method", req.Method),
		zap.String("URI", req.RequestURI),
		zap.Int("Status", respLogger.Status()),
		zap.Int("Size", respLogger.Size()),
	)
}

func newResponseLogger(w http.ResponseWriter) (*responseLogger, http.ResponseWriter) {
	logger := &responseLogger{w: w, status: http.StatusOK}
	return logger, httpsnoop.Wrap(w, httpsnoop.Hooks{
		Write: func(httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return logger.Write
		},
		WriteHeader: func(httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return logger.WriteHeader
		},
	})
}

// Write - io.Writer interface.
func (l *responseLogger) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

// Hijack http.Hijack interface.
func (l *responseLogger) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	conn, rw, err := l.w.(http.Hijacker).Hijack()
	if err == nil && l.status == 0 {
		// The status will be StatusSwitchingProtocols if there was no error and
		// WriteHeader has not been called yet
		l.status = http.StatusSwitchingProtocols
	}
	return conn, rw, err
}
