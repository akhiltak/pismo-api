package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/akhiltak/pismo-api/config"
	_ "github.com/akhiltak/pismo-api/docs"
	"github.com/akhiltak/pismo-api/internal/server"
)

// @title			Transaction API
// @version		1.0
// @description	Transaction API for Pismo
// @contact.name	akhiltak@gmail.com
// @BasePath		/
func main() {
	cfg := config.Get()
	ctx := context.Background()

	if cfg.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	log.Println("Starting Pismo transaction server ...")
	slog.Info("Starting Pismo transaction server ...")

	// run rest server
	s := server.New(ctx, cfg)
	if err := s.Run(cfg.HTTPListenHostPort); err != nil {
		log.Fatal("Error starting server: ", err)
	}

	// wait for shutdown
	waitForShutdown(ctx, s)
}

// waitForShutdown handles graceful termination
func waitForShutdown(ctx context.Context, s *server.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server shut down successfully")
}
