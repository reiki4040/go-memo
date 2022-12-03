package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EnvLogMode         = "SERVER_LOG_MODE"
	EnvLogModeValueDev = "developement"

	EnvLogLevel           = "SERVER_LOG_LEVEL"
	EnvLogLevelValueDebug = "debug"
	EnvLogLevelValueInfo  = "info"
	EnvLogLevelValueWarn  = "warn"
	EnvLogLevelValueError = "error"
)

func main() {
	c := zap.NewProductionConfig()
	c.Development = (strings.ToLower(os.Getenv(EnvLogMode)) == EnvLogModeValueDev)
	switch strings.ToLower(os.Getenv(EnvLogLevel)) {
	case EnvLogLevelValueDebug:
		c.Level.SetLevel(zapcore.DebugLevel)
	case EnvLogLevelValueWarn:
		c.Level.SetLevel(zapcore.WarnLevel)
	case EnvLogLevelValueError:
		c.Level.SetLevel(zapcore.ErrorLevel)
	case EnvLogLevelValueInfo:
		// none (=Info)
	default:
		// none (=Info)
	}
	c.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger, err := c.Build()
	if err != nil {
		panic(fmt.Sprintf("failed initialize logger: %v", err))
	}
	mLogger := logger.Named("main")

	mux, err := NewServeMux(logger.Named("servemux"))
	if err != nil {
		mLogger.Error("failed initialize ServeMux", zap.Error(err))
		return
	}

	s := &http.Server{
		Addr:        ":8080",
		Handler:     mux,
		ReadTimeout: 5 * time.Second,
	}

	// signal handling
	signalCtx, signalCancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer signalCancelFunc()

	// graceful shutdown with timeout.
	waitGraceful := make(chan struct{})
	go func() {
		<-signalCtx.Done()
		logger.Info("got signal. server will do graceful shutdown...")
		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()
		if err := s.Shutdown(ctx); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Warn("failed graceful shutdown with timeout passed",
					zap.Duration("timeout", time.Duration(5*time.Second)),
				)
			} else {
				logger.Error("failed graceful shutdown", zap.Error(err))
			}
		} else {
			logger.Info("OK. Server is shotdown normaly.")
		}
		close(waitGraceful)
	}()

	if err := s.ListenAndServeTLS("./cert.pem", "./cert-key.pem"); err != nil {
		if errors.Is(http.ErrServerClosed, err) {
			mLogger.Warn("got server closed maybe got signal that stop process", zap.Error(err))
		} else {
			mLogger.Error("got error from ListenAndServer() that is NOT server closed", zap.Error(err))
		}
	}

	/*
		got signal and call Shutdown() then ListenAndServe() returns 'Server closed'.
		it mean is stop new requests recieving, however recieved requests still
		processing with until graceful shutdown timeout.
		main() will finish if not waiting Shutdown() finished.
	*/
	<-waitGraceful
}
