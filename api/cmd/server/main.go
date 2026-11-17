package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/nyxsr/meme-portal/api/internal/app"
	"github.com/nyxsr/meme-portal/api/internal/config"
	"github.com/nyxsr/meme-portal/api/internal/platform/logger"
)

func main() {
	err := godotenv.Load()
	if errors.Is(err, os.ErrNotExist) {
		err = godotenv.Load("../.env")
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	lg := logger.Setup(cfg.LogLevel)

	startupContext, startupCancel := context.WithTimeout(context.Background(), 10*time.Second)

	application, err := app.New(startupContext, cfg)
	if err != nil {
		log.Fatalf("Error creating application: %v", err)
	}
	startupCancel()

	serverErrors := make(chan error)

	go func() {
		lg.Info("starting server", "port", cfg.Port)
		serverErrors <- application.Run()
	}()

	signalContext, signalCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer signalCancel()

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
		return
	case <-signalContext.Done():
		lg.Info("shutting down server")
	}

	shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := application.Shutdown(shutdownContext); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
}
