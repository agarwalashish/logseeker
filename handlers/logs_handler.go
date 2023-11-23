package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"logseeker/models"
	"logseeker/services"
	"net/http"

	"github.com/rs/zerolog/log"
)

type LogsHandlerInterface interface {
	SearchRequest(w http.ResponseWriter, r *http.Request)
}

type LogsHandler struct {
	BaseHandler
	SearchService services.SearchServiceInterface
}

func NewLogsHandler(searchService services.SearchServiceInterface) *LogsHandler {
	return &LogsHandler{
		SearchService: searchService,
	}
}

func (lh *LogsHandler) SearchRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Msg("Error reading request body")
		lh.WriteJSON(w, r, errors.New("error reading request body"), http.StatusBadRequest)
		return
	}

	var searchRequest models.SearchRequest
	json.Unmarshal(body, &searchRequest)

	log.Info().Any("searchRequest", searchRequest).Msg("rcvd searchRequest")
	lines, e := lh.SearchService.Search(&searchRequest)
	if e != nil {
		log.Error().Msg(e.Message)
		lh.WriteJSON(w, r, e, e.Code)
		return
	}

	searchResponse := &models.SearchResponse{Lines: lines}
	lh.WriteJSON(w, r, searchResponse, http.StatusOK)
}
