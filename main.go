package main

import (
	"logseeker/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := setupRouter()
	http.ListenAndServe(":3000", router)
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()

	var logsHandler handlers.LogsHandlerInterface = handlers.NewLogsHandler()
	r.Route("/logs", func(r chi.Router) {
		r.Post("/search", logsHandler.SearchRequest)
	})

	return r
}
