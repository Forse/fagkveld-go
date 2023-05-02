package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	// After the archival of gorilla/mux, chi seems to be the spiritual successor
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// This package should run a simple HTTP API for a todo-application.
// The included code is some basic scaffolding for the API,
// with some best practices applied for running in a containerized environment.
// Add the actual APIs to learn how to build HTTP APIs in Go.
// Start in the router function.

func main() {
	// Using builtin logger
	// OpenTelemetry logging not implemented yet for Go
	// other solutions exist, but we prefer simple solutions
	logger := log.New(os.Stdout, "", 0)

	// Read configuration from environment
	cfg, err := configuration()
	if err != nil {
		logger.Fatalf("ERROR: could not read configuration: %s", err)
	}
	logger.Printf("INFO: configuration: %s", cfg)

	// Build HTTP router
	router, err := router(logger)
	if err != nil {
		logger.Fatalf("ERROR: could not create router: %s", err)
	}

	// Boot HTTP server
	logger.Printf("INFO: serving HTTP...")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: router,
	}
	go func(server *http.Server) {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("ERROR: while starting server: %s", err)
		}
	}(server)

	// Handle OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)    // Ctrl-c
	signal.Notify(quit, syscall.SIGTERM) // Kubernetes initially sends SIGTERM
	// After SIGTERM, we may get SIGKILL, but at that point the OS will have killed the process

	// After receiving the OS signal, attempt graceful shutdown with a timeout
	// timeout may reflect 'terminationGracePeriodSeconds' in a Kubernetes environment (default is 30s)
	sig := <-quit
	logger.Printf("INFO: received OS signal (will shutdown): %s", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("ERROR: error shutting down server: %s", err)
	}
	logger.Printf("INFO: shutdown gracefully")
}

func router(logger *log.Logger) (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy"))
	})

	err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logger.Printf("INFO: route: %s %s\n", method, route)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r, nil
}

type config struct {
	port int
}

// The String function will be called when string formatting is used
func (c *config) String() string {
	return fmt.Sprintf("port=%d", c.port)
}

// A simple function accepting configuration from environment variables
// an alternative solution would implement this: https://github.com/gookit/config
// but in Go we prefer simple solutions
func configuration() (*config, error) {
	var port int
	portStr := os.Getenv("PORT")
	if portStr == "" {
		port = 3000
	} else {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("ERROR: invalid port config: %w", err)
		}

		port = p
	}

	return &config{port: port}, nil
}
