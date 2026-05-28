package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/findmesh/findmesh/backend/internal/config"
	"github.com/findmesh/findmesh/backend/internal/db"
	"github.com/findmesh/findmesh/backend/internal/sightings"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := db.NewMemoryStore()
	svc := sightings.NewService(store)
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	logger.Info("findmesh retention worker started")
	for range ticker.C {
		deleted := svc.Retain(cfg.RawRetention, cfg.LostRetention)
		logger.Info("retention cycle complete", "deleted_raw_sightings", deleted)
	}
}
