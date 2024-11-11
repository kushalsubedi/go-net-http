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

	"github.com/kushalsubedi/students-api/internal/config"
	"github.com/kushalsubedi/students-api/internal/http/handlers/student"
)

func main() {

	// load config
	cfg := config.MustLoad()
	// database setup
	// setup router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New())
	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	slog.Info("server started", slog.String("address", cfg.Addr))
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server ")
		}
	}()

	<-done
	slog.Info("Shutting Down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to Shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown succesfully")
}
