package main

import (
	"logseeker/handlers"
	"logseeker/services"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
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
