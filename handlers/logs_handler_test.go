package handlers

import (
	"bytes"
	"logseeker/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSearchService struct {
}

func (m *mockSearchService) Search(req *models.SearchRequest) ([]string, *models.Error) {
	lines := []string{
		"Line 1",
		"Line 2",
	}

	return lines, nil
}

func TestSearchRequest(t *testing.T) {
	searchService := &mockSearchService{}

	logsHandler := NewLogsHandler(searchService)
	testCases := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "valid request",
			requestBody:    validRequestBody,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest("POST", "/search", bytes.NewBufferString(tc.requestBody))
		rr := httptest.NewRecorder()

		logsHandler.SearchRequest(rr, req)
		assert.Equal(t, tc.expectedStatus, rr.Code)
	}
}

var validRequestBody string = `{
	"numLines": 500,
	"filename": "../tests/data/log_file.log",
	"keywords": "user"
}`
