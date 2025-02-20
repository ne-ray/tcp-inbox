// Package app configures and runs application.
package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ne-ray/tcp-inbox/config"
	// "github.com/ne-ray/tcp-inbox/internal/usecase"
	// "github.com/evrone/go-clean-template/internal/usecase/repo"
	// "github.com/evrone/go-clean-template/internal/usecase/webapi"
	// "github.com/evrone/go-clean-template/pkg/httpserver"
	"github.com/ne-ray/tcp-inbox/pkg/logger"
)

const (
	APP_NAME    = "TCP inbox"
	APP_VERSION = "1.0.0"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level, APP_NAME, APP_VERSION)

	l.Debug("Running application...")

	// Use case
	// translationUseCase := usecase.New(
	// 	repo.New(pg),
	// 	webapi.New(),
	// )

	// TCP Server
	// handler := gin.New()
	// v1.NewRouter(handler, l, translationUseCase)
	// httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	l.Info("Application started")

	var err error
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
		// case err = <-httpServer.Notify():
		// 	l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	// err = httpServer.Shutdown()
	// if err != nil {
	// 	l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	// }

	if err = l.Shutdown(); err != nil {
		log.Fatalf("Logger shutdown error: %s", err)
	}
}
