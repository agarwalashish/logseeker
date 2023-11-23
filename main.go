package main

import (
	"logseeker/handlers"
	"logseeker/services"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	setupLogger()
	router := setupRouter()
	http.ListenAndServe(":8080", router)
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Set to the frontend's URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(cors.Handler)

	// Initialize the services
	searchService := services.NewSearchService()

	// Initialize the handlers
	var logsHandler handlers.LogsHandlerInterface = handlers.NewLogsHandler(searchService)

	// Initialize the routes
	r.Route("/logs", func(r chi.Router) {
		r.Post("/search", logsHandler.SearchRequest)
	})

	return r
}

func setupLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
