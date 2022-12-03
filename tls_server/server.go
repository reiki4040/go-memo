package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func NewServeMux(l *zap.Logger) (*http.ServeMux, error) {
	logging, err := NewLoggingWith(l)
	if err != nil {
		l.Error("failed logging middleware", zap.Error(err))
	}

	mux := http.NewServeMux()
	api := &API{}
	mux.Handle("/", logging.RequestLogging(http.HandlerFunc(api.HelloWorld)))
	mux.Handle("/sleep15", logging.RequestLogging(http.HandlerFunc(api.Sleep)))

	return mux, nil
}

type ReqContextLogger struct{}

func NewLoggingWith(l *zap.Logger) (*Logging, error) {
	return &Logging{
		logger: l,
	}, nil
}

type Logging struct {
	logger *zap.Logger
}

func (l *Logging) RequestLogging(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), ReqContextLogger{}, l.logger)
		start := time.Now()
		h.ServeHTTP(w, req.WithContext(ctx))
		responseTime := time.Since(start)

		l.logger.Info("request log",
			zap.Time("request_time", start),
			zap.String("remote_addr", req.RemoteAddr),
			zap.String("host", req.Host),
			zap.String("method", req.Method),
			zap.String("path", req.RequestURI),
			zap.String("referer", req.Referer()),
			zap.String("user_agent", req.UserAgent()),
			zap.Int64("request_size", req.ContentLength),
			zap.Any("form", req.Form),
			zap.Duration("response_time", responseTime),
		)
	}

	return http.HandlerFunc(fn)
}

func GetLoggerFromContext(ctx context.Context) (*zap.Logger, error) {
	lv := ctx.Value(ReqContextLogger{})
	l, ok := lv.(*zap.Logger)
	if !ok {
		return nil, fmt.Errorf("does not set logger in request context, please chack middleware")
	}

	return l, nil
}

type API struct{}

func (api *API) HelloWorld(w http.ResponseWriter, req *http.Request) {
	l, err := GetLoggerFromContext(req.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Hello world"))
	l.Debug("called hello world")
}

func (api *API) Sleep(w http.ResponseWriter, req *http.Request) {
	time.Sleep(15 * time.Second)
	w.Write([]byte("wake up. I was sleeping until 15 seconds."))
}
