package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gurveer1510/student-api/internal/config"
	"github.com/Gurveer1510/student-api/internal/http/handlers/student"
	"github.com/Gurveer1510/student-api/internal/storage/sqlite"
)

func main() {
	// Load config
	cfg := config.MustLoad()
	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialiazed", slog.String("env",cfg.Env))
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	// setup server
	server := http.Server {
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("server started ---- >>>", slog.String("address", cfg.Addr))

	// GRACEFUL SHUTDOWN
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	go func () {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed{
			log.Fatalf("Failed to start the server: %s", err.Error())
		}
	}()

	<-done
	slog.Info("shutting down the server")

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	err = server.Shutdown(ctxWithTimeout)
	if err != nil {
		slog.Error("Failed to shutdown the server", slog.String("error", err.Error()))
	}
}	