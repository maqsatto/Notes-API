package main

import (
	"context"
	"fmt"
	"log"

	"github.com/maqsatto/Notes-API/internal/config"
	"github.com/maqsatto/Notes-API/internal/database"
)


func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("DB_NAME =", cfg.Database.DBName)
	log.Println("DSN =", cfg.Database.ConnectionString())
	ctx := context.Background()
	//DB Connect
	db, err := database.NewPostgresDB(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("DB connected")
}

