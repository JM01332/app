package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/JM01332/app/internal/app"
	"github.com/JM01332/app/internal/config"
	"go.uber.org/zap"
)

const readHeaderTimeout = 5 * time.Second

const startupBanner = `
   ______                _              ___    ____  ____
  / ____/___ ___________(_)__  _____   /   |  / __ \/  _/
 / /   / __ ` + "`" + `/ ___/ ___/ / _ \/ ___/  / /| | / /_/ // /
/ /___/ /_/ / /  / /  / /  __/ /     / ___ |/ ____// /
\____/\__,_/_/  /_/  /_/\___/_/     /_/  |_/_/   /___/
`

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

	printBanner(serverConfig.Port)
	logger.Info("starting API server", zap.String("address", server.Addr))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("API server stopped", zap.Error(err))
	}
}

func printBanner(port string) {
	baseURL := "http://" + net.JoinHostPort("localhost", port)
	fmt.Fprintf(
		os.Stdout,
		"%s\nCarrier API läuft auf %s\nHealth: %s/health\n\n",
		startupBanner,
		baseURL,
		baseURL,
	)
}
