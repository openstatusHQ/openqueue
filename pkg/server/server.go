package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/openstatushq/openqueue/pkg/database"
	"github.com/rs/zerolog/log"
)

type Server struct {
	port int
	dbs map[string]*sqlx.DB
}

type Options struct {
	Port int
	Queues []struct{
		Name string
		DB string
	}
}
func NewServer(opts Options) error {


	s := new(Server)
	s.port = opts.Port
	for _, q := range opts.Queues {
		db := database.GetDatabase(q.DB)
		if db == nil {
			log.Fatal().Msgf("Error setting up database %s", q.DB)
		}
		s.dbs[q.Name] = db
	}

	server := newServer(s)

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal().Msg("graceful shutdown timed out.. forcing exit.")

			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal().Err(err).Msg("server shutdown failed")
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server failed")
		return err
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	return nil
}

func newServer(s *Server) *http.Server {


	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
