package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/floqast/task-management/backend/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	application, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	application.Logger.Info("server starting", "port", port)

	handler := app.SetupRouter(application, port)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		application.Logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
