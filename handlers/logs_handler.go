package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"voltron/models"
	"voltron/services"
)

type LogsHandler struct {
	BaseHandler
	SearchService services.SearchServiceInterface
}

func NewLogsHandler() *LogsHandler {
	searchService := services.NewSearchService()
	return &LogsHandler{
		SearchService: searchService,
	}
}

func (lh *LogsHandler) SearchRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		//TODO: handle error
	}

	var searchRequest models.SearchRequest
	err = json.Unmarshal(body, &searchRequest)

	log.Printf("Received: %+v\n", searchRequest)
	lines, err := lh.SearchService.Search(&searchRequest)
	if err != nil {
		lh.SendJSON(w, r, err, http.StatusBadRequest)
		return
	}

	lh.SendJSON(w, r, lines, http.StatusOK)
}
