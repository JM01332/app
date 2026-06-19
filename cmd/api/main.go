package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/JM01332/app/internal/app"
	carrierservice "github.com/JM01332/app/internal/carrier/service"
	"github.com/JM01332/app/internal/config"
	"github.com/JM01332/app/internal/database"
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

	appConfig, err := config.Load()
	if err != nil {
		logger.Fatal("load application configuration", zap.Error(err))
	}

	postgres, err := database.OpenPostgres(context.Background(), appConfig.DatabaseURL)
	if err != nil {
		logger.Fatal("connect to PostgreSQL", zap.Error(err))
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			logger.Error("close PostgreSQL connection", zap.Error(err))
		}
	}()

	carrierRepository := carrierservice.NewCarrierRepository(postgres.DB)
	carrierService := carrierservice.NewCarrierService(carrierRepository)

	server := &http.Server{
		Addr:              net.JoinHostPort("", appConfig.Port),
		Handler:           app.NewRouter(carrierService),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	printBanner(appConfig.Port)
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
