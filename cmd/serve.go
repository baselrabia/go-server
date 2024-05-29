package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baselrabia/go-server/config"
	"github.com/baselrabia/go-server/handlers"
	"github.com/baselrabia/go-server/internal/persistence"
	"github.com/baselrabia/go-server/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	persistor, err := persistence.NewPersistor(cfg.DataFile)
	if err != nil {
		log.Fatalf("Failed to create persistor: %v", err)
	}
	defer persistor.Close()

	srv, err := server.NewServer(cfg.WindowDuration, cfg.PersistInterval, persistor)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	http.HandleFunc("/", handlers.CounterHandler(srv))

	log.Printf("Starting server on port %s", cfg.Port)
	httpSrv := &http.Server{
		Addr: cfg.Port,
	}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Handle shutdown signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Persist data before shutting down
	srv.PersistData()

	// Create a context with a timeout to give the server a chance to
    // finish up any ongoing requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shut down server: %v", err)
	}

	log.Println("Server stopped")
}
