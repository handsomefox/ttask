package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/handsomefox/ttask/api/handler"
	"github.com/handsomefox/ttask/api/middleware"
	"github.com/julienschmidt/httprouter"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	slog.SetDefault(log)

	// Use https://github.com/julienschmidt/httprouter for creating a server
	r := httprouter.New()

	// 1. Create a REST endpoint called calculate
	r.POST("/calculate", middleware.Calculate(handler.Calculate))

	server := &http.Server{
		Addr:              ":8989", // available at port 8989 so we can access it http://localhost:8989/calculate
		Handler:           r,
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       5 * time.Second,
	}

	slog.Info("Starting server", "addr", "localhost:8989")

	if err := server.ListenAndServe(); err != nil {
		slog.Error("fatal server error", "err", err)
		os.Exit(1)
	}
}
