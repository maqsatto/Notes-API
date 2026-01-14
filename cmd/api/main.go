package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/maqsatto/Notes-API/internal/config"
	"github.com/maqsatto/Notes-API/internal/database"
	"github.com/maqsatto/Notes-API/internal/http/router"
	"github.com/maqsatto/Notes-API/internal/logger"
)


func main() {
	//config
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// app logger
	logg, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer logg.Close()

	//DB Connect
	ctx := context.Background()
	db, err := database.NewPostgresDB(ctx, cfg.Database)
	if err != nil {
		logg.Error("failed to connect to database", err)
		return
	}
	defer db.Close()
	fmt.Println("DB connected")

	// build Server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	h := router.New(router.Deps{
		Config: cfg,
		Logger: logg,
		DB: db,
	})
	srv := &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	//graceful shutdown

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logg.Info("server started on " + addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Error("server error", err)
			stop() // trigger shutdown if server failed unexpectedly
		}
	}()
	<-ctx.Done()
	logg.Info("shutdown signal received")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logg.Error("server shutdown failed", err)
		srv.Close()
	}
	logg.Info("server exited")
}
