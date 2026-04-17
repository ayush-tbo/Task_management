package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/floqast/task-management/backend/internal/app"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	application, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	application.Logger.Printf("We are running on port %s \n", port)

	r := app.SetupRoutes(application)

	r.Get("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message":"Hello from Chi backend!"}`)
	})

	r.Get("/swagger/doc.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		http.ServeFile(w, r, "docs/swagger.yaml")
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.yaml"),
	))
	application.Logger.Printf("Swagger UI: http://localhost:%s/swagger/", port)

	corsHandler := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      corsHandler(r),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		application.Logger.Fatal(err)
	}
}
