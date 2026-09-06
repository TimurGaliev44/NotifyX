package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TimurGaliev44/notifyx/internal/ingestion/handler"
	"github.com/TimurGaliev44/notifyx/internal/ingestion/producer"
	"github.com/TimurGaliev44/notifyx/internal/ingestion/repository"
	"github.com/TimurGaliev44/notifyx/internal/ingestion/service"
	"github.com/TimurGaliev44/notifyx/internal/pkg/config"
	"github.com/TimurGaliev44/notifyx/internal/pkg/kafka"
	"github.com/TimurGaliev44/notifyx/internal/pkg/logger"
	"github.com/TimurGaliev44/notifyx/internal/pkg/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log := logger.New("ingestion-api", slog.LevelInfo)

	startCtx, startCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startCancel()

	db, err := postgres.New(startCtx, cfg.PostgresDSN)
	if err != nil {
		log.Error("failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	kafkaClient, err := kafka.New(cfg.KafkaAddrs)
	if err != nil {
		log.Error("failed to create kafka client: %v", err)
		os.Exit(1)
	}

	prod, err := kafkaClient.NewProducer()
	if err != nil {
		log.Error("failed to create kafka producer: %v", err)
		os.Exit(1)
	}

	repo := repository.New(db)
	pub := producer.New(prod)

	svc := service.New(repo, pub)
	h := handler.New(svc)
	server := http.Server{
		Addr:              cfg.HTTPaddr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info("server closed")
				os.Exit(0)
			}
			log.Error("failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	server.Shutdown(shutdownCtx)
	prod.Close()
	kafkaClient.Close()
	db.Close()
}
