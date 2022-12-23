package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"workshop/cmd/server/config"
	"workshop/internal/storage"
	"workshop/pkg/dbcollector"
	"workshop/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"workshop/internal/handlers"
	"workshop/internal/users"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	defaultLog := logger.DefaultLogger()

	var cfg, help, err = config.New()
	if err != nil {
		if help != "" {
			defaultLog.Fatal().Msg(help.String())
		}
		defaultLog.Fatal().Err(err).Msg("failed to parse config")
	}

	log, err := logger.New(cfg.Log)
	if err != nil {
		defaultLog.Fatal().Err(err).Msg("failed to init logger")
	}

	log.Info().Msg("starting")
	defer log.Info().Msg("shutdown")

	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer func() {
		_ = db.Close()
	}()

	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)

	db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdleTime)

	ur := storage.NewStorage(db)
	us := users.NewService(ur)
	uh := handlers.NewUsers(us)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(logger.LoggerMiddleware(log))

	prometheus.MustRegister(dbcollector.NewSQLDatabaseCollector("general", "main", "sqlite", db))
	r.Mount("/metrics", promhttp.Handler())
	r.Mount("/users", uh.Routes())

	s := http.Server{
		Addr:    cfg.HTTP.Addr,
		Handler: r,

		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,

		IdleTimeout: cfg.HTTP.IdleTimeout,
	}

	go func() {
		fmt.Println(s.ListenAndServe())
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	<-ctx.Done()
	log.Info().Msg("signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefullTimeout)
	defer cancel()

	log.Info().Msg("shutting down")
	_ = s.Shutdown(ctx)
}
