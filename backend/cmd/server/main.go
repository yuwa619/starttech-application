package main

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/starttech/starttech-application/backend/internal/config"
	"github.com/starttech/starttech-application/backend/internal/handlers"
	"github.com/starttech/starttech-application/backend/internal/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mongoClient := connectMongo(ctx, cfg, logger)
	if mongoClient != nil {
		defer mongoClient.Disconnect(context.Background())
	}

	redisClient := connectRedis(ctx, cfg, logger)
	if redisClient != nil {
		defer redisClient.Close()
	}

	router := handlers.NewRouter(handlers.Dependencies{
		Logger:      logger,
		MongoClient: mongoClient,
		RedisClient: redisClient,
		Environment: cfg.Environment,
		ServiceName: cfg.ServiceName,
		StartedAt:   time.Now(),
	})

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           logging.Middleware(logger, router),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("server starting", "port", cfg.Port, "environment", cfg.Environment)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown requested")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}
}

func connectMongo(ctx context.Context, cfg config.Config, logger *slog.Logger) *mongo.Client {
	if cfg.MongoDBURI == "" {
		logger.Warn("mongodb uri not configured")
		return nil
	}

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(cfg.MongoDBURI))
	if err != nil {
		logger.Error("mongodb connection failed", "error", err)
		return nil
	}

	return client
}

func connectRedis(ctx context.Context, cfg config.Config, logger *slog.Logger) *redis.Client {
	if cfg.RedisAddr == "" {
		logger.Warn("redis address not configured")
		return nil
	}

	options := &redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}

	if cfg.RedisTLS {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client := redis.NewClient(options)

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx).Err(); err != nil {
		logger.Error("redis connection failed", "error", err)
		return nil
	}

	return client
}
