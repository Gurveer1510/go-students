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
)

func main() {
	// Load config
	cfg := config.MustLoad()
	// database setup
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Hello from the student-server"))
	})
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

	err := server.Shutdown(ctxWithTimeout)
	if err != nil {
		slog.Error("Failed to shutdown the server", slog.String("error", err.Error()))
	}
}	