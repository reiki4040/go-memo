package main

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewServeMux() (*http.ServeMux, error) {
	mux := http.NewServeMux()
	mux.Handle("/", reqLog(middle(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello world"))
	}))))
	mux.Handle("/sleep15", reqLog(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(15 * time.Second)
		w.Write([]byte("wake up. I was sleeping until 15 seconds."))
	})))

	return mux, nil
}

func middle(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("authorization")
		// とりあえずヘッダの有無だけ
		if auth == "" {
			http.Error(w, "Unauthorize request", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

type ReqContextLogger struct{}

func reqLog(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		c := zap.NewProductionConfig()
		c.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
		l, _ := c.Build()

		_ = context.WithValue(req.Context(), ReqContextLogger{}, l)
		start := time.Now()
		h.ServeHTTP(w, req)
		responseTime := time.Since(start)

		l.Info("request log",
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
