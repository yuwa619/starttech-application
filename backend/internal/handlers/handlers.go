package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dependencies struct {
	Logger      *slog.Logger
	MongoClient *mongo.Client
	RedisClient *redis.Client
	Environment string
	ServiceName string
	StartedAt   time.Time
}

func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", deps.health)
	mux.HandleFunc("GET /readyz", deps.ready)
	mux.HandleFunc("GET /api/v1/status", deps.status)
	return mux
}

func (d Dependencies) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":         "ok",
		"service":        d.ServiceName,
		"environment":    d.Environment,
		"uptime_seconds": int(time.Since(d.StartedAt).Seconds()),
	})
}

func (d Dependencies) status(w http.ResponseWriter, r *http.Request) {
	d.health(w, r)
}

func (d Dependencies) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	checks := map[string]string{}
	statusCode := http.StatusOK

	if d.MongoClient != nil {
		if err := d.MongoClient.Ping(ctx, nil); err != nil {
			checks["mongodb"] = "unavailable"
			statusCode = http.StatusServiceUnavailable
			d.Logger.Error("mongodb readiness check failed", "error", err)
		} else {
			checks["mongodb"] = "ok"
		}
	} else {
		checks["mongodb"] = "not_configured"
	}

	if d.RedisClient != nil {
		if err := d.RedisClient.Ping(ctx).Err(); err != nil {
			checks["redis"] = "unavailable"
			statusCode = http.StatusServiceUnavailable
			d.Logger.Error("redis readiness check failed", "error", err)
		} else {
			checks["redis"] = "ok"
		}
	} else {
		checks["redis"] = "not_configured"
	}

	writeJSON(w, statusCode, map[string]any{
		"status": statusFromCode(statusCode),
		"checks": checks,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func statusFromCode(statusCode int) string {
	if statusCode >= http.StatusBadRequest {
		return "degraded"
	}
	return "ready"
}
