package main

import (
	"context"
	"fmt"
	"log"

	"github.com/maqsatto/Notes-API/internal/config"
	"github.com/maqsatto/Notes-API/internal/database"
	"github.com/maqsatto/Notes-API/internal/migration"
)


func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()
	//DB Connect
	db, err := database.NewPostgresDB(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// MIGRATE UP
	if err := migration.MigrateUp(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("DB connected + migrations applied")
}

