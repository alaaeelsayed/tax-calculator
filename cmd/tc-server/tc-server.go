package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alaaeelsayed/tax-calculator/internal/api"
	"github.com/alaaeelsayed/tax-calculator/internal/client"
	"github.com/alaaeelsayed/tax-calculator/internal/config"
	"github.com/alaaeelsayed/tax-calculator/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	calculator := service.NewTaxCalculatorService(client.NewClient(cfg.TaxAPIURL))
	apiServer := api.NewServer(calculator)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: apiServer.SetupRoutes(),
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
