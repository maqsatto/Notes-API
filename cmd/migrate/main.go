package main

import (
	"context"
	"flag"
	"log"

	"github.com/maqsatto/Notes-API/internal/config"
	"github.com/maqsatto/Notes-API/internal/database"
	"github.com/maqsatto/Notes-API/internal/migration"
)

func main() {
	// usage: go run ./cmd/migrate -direction=up
	// usage: go run ./cmd/migrate -direction=down
	direction := flag.String("direction", "up", "migration direction: up or down")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	ctx := context.Background()
	db, err := database.NewPostgresDB(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	switch *direction {
	case "up":
		if err := migration.MigrateUp(db); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migrations UP completed")
	case "down":
		if err := migration.MigrateDown(db); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migrations DOWN completed (one step)")
	default:
		log.Fatalf("Unknown direction: %s (use up|down)", *direction)
	}
}
