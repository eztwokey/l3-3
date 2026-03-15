package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eztwokey/l3-3/internal/api"
	"github.com/eztwokey/l3-3/internal/config"
	"github.com/eztwokey/l3-3/internal/logic"
	"github.com/eztwokey/l3-3/internal/storage"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/logger"
)

func main() {
	cfg := new(config.Config)
	if err := cfg.Read(config.LocalPath); err != nil {
		log.Fatal(err)
	}

	wbLog := logger.NewSlogAdapter("l3-commenttree", "local")

	db, err := dbpg.New(cfg.Postgres.DSN(), nil, nil)
	if err != nil {
		log.Fatal("postgres connect:", err)
	}
	defer func() {
		if err := db.Master.Close(); err != nil {
			wbLog.Error("postgres close failed", "err", err)
		}
	}()

	if err := db.Master.Ping(); err != nil {
		log.Fatal("postgres ping:", err)
	}
	wbLog.Info("postgres connected")

	if err := runMigrations(db); err != nil {
		log.Fatal("migrations failed:", err)
	}
	wbLog.Info("migrations applied")

	store := storage.New(db)
	logic := logic.New(store, wbLog)
	http := api.New(cfg, logic, wbLog)

	errChan := make(chan error, 1)
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		errChan <- http.Run()
	}()

	wbLog.Info("server started", "addr", cfg.Api.Addr)

	select {
	case sig := <-termChan:
		wbLog.Warn("got term signal", "signal", sig)
	case err := <-errChan:
		wbLog.Warn("got error", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- http.Shutdown(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Fatalf("server forced to shutdown: %v", err)
		}
		wbLog.Info("server shutdown gracefully")
	case <-ctx.Done():
		wbLog.Info("shutdown timeout exceeded")
	}
}

func runMigrations(db *dbpg.DB) error {
	sqlBytes, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(context.Background(), string(sqlBytes))
	return err
}
