package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"workshop/cmd/server/config"
	"workshop/internal/storage"
	"workshop/internal/transport/grpc/pb"
	"workshop/internal/transport/grpc/server"
	"workshop/internal/transport/http/handlers"
	"workshop/pkg/dbcollector"
	"workshop/pkg/logger"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	uh := handlers.NewUsers(us, ur)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(logger.LoggerMiddleware(log))

	prometheus.MustRegister(dbcollector.NewSQLDatabaseCollector("general", "main", "sqlite", db))
	r.Mount("/metrics", promhttp.Handler())

	r.Mount("/users", uh.Routes())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	httpServer := http.Server{
		Addr:    cfg.HTTP.Addr,
		Handler: r,

		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,

		IdleTimeout: cfg.HTTP.IdleTimeout,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Err(err).Msg("http server error received")
			stop()
		}
	}()

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	pb.RegisterUsersServer(grpcServer, server.NewUsers(us, ur))

	go func() {
		lis, err := net.Listen("tcp", cfg.GRPC.Addr)
		if err != nil {
			log.Err(err).Msg("failed to start tcp server")
			stop()
		}

		if err := grpcServer.Serve(lis); err != grpc.ErrServerStopped {
			log.Err(err).Msg("failed to start grpc server")
			stop()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("signal received")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefullTimeout)
	defer cancel()

	log.Info().Msg("shutting down")
	_ = httpServer.Shutdown(ctx)
	grpcServer.GracefulStop()
}
