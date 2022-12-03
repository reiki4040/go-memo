package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewServeMux() (*http.ServeMux, error) {
	mux := http.NewServeMux()
	mux.Handle("/", reqLog(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		l, err := getLoggerFromContext(req.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		l.Debug("call /")
		w.Write([]byte("Hello world"))
	})))
	mux.Handle("/sleep15", reqLog(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		time.Sleep(15 * time.Second)
		w.Write([]byte("wake up. I was sleeping until 15 seconds."))
	})))

	return mux, nil
}

type ReqContextLogger struct{}

func reqLog(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		c := zap.NewProductionConfig()
		c.Level.SetLevel(zapcore.DebugLevel)
		c.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
		l, _ := c.Build()

		ctx := context.WithValue(req.Context(), ReqContextLogger{}, l)
		start := time.Now()
		h.ServeHTTP(w, req.WithContext(ctx))
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

func getLoggerFromContext(ctx context.Context) (*zap.Logger, error) {
	lv := ctx.Value(ReqContextLogger{})
	l, ok := lv.(*zap.Logger)
	if !ok {
		return nil, fmt.Errorf("does not set logger in request context, please chack middleware")
	}

	return l, nil
}
