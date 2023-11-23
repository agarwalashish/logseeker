package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"logseeker/models"
	"logseeker/services"
	"net/http"

	"go.uber.org/zap"
)

type LogsHandlerInterface interface {
	SearchRequest(w http.ResponseWriter, r *http.Request)
}

type LogsHandler struct {
	BaseHandler
	SearchService services.SearchServiceInterface
	logger        *zap.Logger
}

func NewLogsHandler(logger *zap.Logger) *LogsHandler {
	searchService := services.NewSearchService(logger)
	return &LogsHandler{
		SearchService: searchService,
		logger:        logger,
	}
}

func (lh *LogsHandler) SearchRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		lh.logger.Error("Error reading request body")
		lh.WriteError(w, r, errors.New("Error reading request body"), http.StatusBadRequest)
		return
	}

	var searchRequest models.SearchRequest
	err = json.Unmarshal(body, &searchRequest)

	lh.logger.Info("Rcvd searchRequest", zap.Any("searchRequest", searchRequest))
	lines, e := lh.SearchService.Search(&searchRequest)
	if e != nil {
		lh.logger.Error(e.Message)
		lh.WriteJSON(w, r, e, e.Code)
		return
	}

	searchResponse := &models.SearchResponse{Lines: lines}
	lh.WriteJSON(w, r, searchResponse, http.StatusOK)
}
