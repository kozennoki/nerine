package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kozennoki/nerine/internal/infrastructure/config"
	"github.com/kozennoki/nerine/internal/infrastructure/logger"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	zapLogger, err := logger.New()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer zapLogger.Sync()

	cfg, err := config.Load()
	if err != nil {
		zapLogger.Fatal("Failed to load config", zap.Error(err))
	}

	e := setupServer(cfg, zapLogger)

	if err := startServer(e, cfg, zapLogger); err != nil {
		zapLogger.Fatal("Failed to start server", zap.Error(err))
	}
}
