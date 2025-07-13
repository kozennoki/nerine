package main

import (
	"github.com/kozennoki/nerine/internal/infrastructure/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func setupServer(cfg *config.Config, logger *zap.Logger) *echo.Echo {
	e := echo.New()

	di := NewDIContainer(cfg)
	setupRoutes(e, di, cfg)

	return e
}

func startServer(e *echo.Echo, cfg *config.Config, logger *zap.Logger) error {
	logger.Info("Starting server", zap.String("port", cfg.Port))
	return e.Start(":" + cfg.Port)
}
