// Package app configures and runs application.
package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ne-ray/tcp-inbox/config"
	"github.com/ne-ray/tcp-inbox/internal/usecase/word-of-wisdom"
	"github.com/ne-ray/tcp-inbox/internal/repo/persistent"
	v1 "github.com/ne-ray/tcp-inbox/internal/controller/tcp/v1"
	"github.com/ne-ray/tcp-inbox/pkg/logger"
	"github.com/ne-ray/tcp-inbox/pkg/tcpserver"
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
	wordOfWisdomUseCase := wordofwisdom.New(
		persistent.New(cfg.Storage),
		l,
	)

	// TCP Handler
	h := v1.NewWordOfWisdomHandler(wordOfWisdomUseCase, l, cfg.Session)

	// TCP Server
	tcpserver := tcpserver.New(h, tcpserver.Host(cfg.Host), tcpserver.Port(cfg.Port))
	l.Debug("Start listen host: " + cfg.Host + " port: " + cfg.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	l.Info("Application started")

	var err error
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-tcpserver.Notify():
		l.With("error", err).Error("app - Run tcpServer.Notify error")
	}

	// Shutdown
	err = tcpserver.Shutdown()
	if err != nil {
		l.With("error", err).Error("app - Run - tcpServer.Shutdown error")
	}

	if err = l.Shutdown(); err != nil {
		log.Fatalf("Logger shutdown error: %s", err)
	}

	l.Info("Application stoped")
}
