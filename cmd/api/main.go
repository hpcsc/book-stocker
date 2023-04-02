package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/go-playground/validator/v10"
	"github.com/hpcsc/book-stocker/internal/config"
	"github.com/hpcsc/book-stocker/internal/info"
	"github.com/hpcsc/book-stocker/internal/stock"
	"github.com/hpcsc/book-stocker/internal/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := log.With().Str("version", info.Version).Logger()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:              addr,
		Handler:           serverHandler(cfg, logger),
		ReadHeaderTimeout: 60 * time.Second,
	}

	withCancelCtx, cancelServer := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		s := <-sig
		logger.Info().Msgf("received %v signal", s)

		// Shutdown signal with grace period of 30 seconds
		withTimeoutCtx, cancelTimeout := context.WithTimeout(withCancelCtx, 30*time.Second)
		defer cancelTimeout()

		go func() {
			<-withTimeoutCtx.Done()
			if withTimeoutCtx.Err() == context.DeadlineExceeded {
				logger.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(withTimeoutCtx); err != nil {
			logger.Fatal().Msgf("failed to gracefully shutdown server: %v", err)
		}
		cancelServer()
	}()

	// Run the server
	logger.Info().Msgf("listening at %v", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err)
	}

	// Wait for server context to be stopped
	<-withCancelCtx.Done()
}

func serverHandler(cfg *config.Configuration, logger zerolog.Logger) http.Handler {
	r := chi.NewRouter()
	v := validator.New()

	r.Use(httplog.RequestLogger(httplog.NewLogger("api", httplog.Options{
		JSON:    true,
		Concise: true,
	})))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	info.RegisterRoutes(r)

	s, err := store.NewDynamoDbStore(cfg)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	stock.RegisterRoutes(r, v, s)

	return r
}
