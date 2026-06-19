package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/JM01332/app/internal/app"
	"github.com/JM01332/app/internal/config"
	"go.uber.org/zap"
)

const readHeaderTimeout = 5 * time.Second

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	serverConfig, err := config.LoadServer()
	if err != nil {
		logger.Fatal("load server configuration", zap.Error(err))
	}

	server := &http.Server{
		Addr:              net.JoinHostPort("", serverConfig.Port),
		Handler:           app.NewRouter(),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	logger.Info("starting API server", zap.String("address", server.Addr))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("API server stopped", zap.Error(err))
	}
}
