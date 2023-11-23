package services

import (
	"io"
	"logseeker/models"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

type SearchServiceInterface interface {
	Search(request *models.SearchRequest) ([]string, *models.Error)
}

type SearchService struct {
}

func NewSearchService() *SearchService {
	return &SearchService{}
}

const (
	defaultLineCount = 10
	chunkSize        = 1024
)

func (ss *SearchService) Search(request *models.SearchRequest) ([]string, *models.Error) {
	if request == nil || request.Filename == "" {
		log.Error().Msg("filename is missing from request")
		return nil, &models.Error{Message: "filename is missing", Code: http.StatusBadRequest}
	}

	lineCount := defaultLineCount
	if request.NumLines > 0 {
		lineCount = request.NumLines
	}

	return SearchFile(request.Filename, lineCount, request.Keywords)
}

// SearchFile searches for keywords in the lines from a file
func SearchFile(filename string, lineCount int, phrase string) ([]string, *models.Error) {
	// Check if file exists before opening it
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, &models.Error{Message: "file does not exist", Code: http.StatusBadRequest}
	}

	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, &models.Error{Message: "error opening file", Code: http.StatusInternalServerError}
	}
	defer file.Close()

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return nil, &models.Error{Message: "file is empty", Code: http.StatusBadRequest}
	}

	var lines []string
	var partialLine string
	var currentPos int64 = fileSize

	for currentPos > 0 {
		var startPos int64
		var sizeToRead int64 = chunkSize

		// Adjust the chunk size if we are near the start of the file
		if currentPos < chunkSize {
			sizeToRead = currentPos
		}
		startPos = currentPos - sizeToRead

		// Create a buffer to hold the data read from the file
		buffer := make([]byte, sizeToRead)
		_, err := file.Seek(startPos, io.SeekStart)
		if err != nil {
			return nil, &models.Error{Message: err.Error(), Code: http.StatusInternalServerError}
		}
		_, err = file.Read(buffer)
		if err != nil {
			return nil, &models.Error{Message: err.Error(), Code: http.StatusInternalServerError}
		}

		// Update the current position in the file
		currentPos -= sizeToRead
		buffer = append(buffer, partialLine...)

		for i := len(buffer) - 1; i >= 0; i-- {
			// Check for newline characters to identify lines
			if buffer[i] == '\n' {
				line := string(buffer[i+1:])
				if len(line) > 0 {
					// Check if the keywords are present in the line
					if CheckForKeywords(line, phrase) {
						// Prepend the line to the lines slice
						lines = append(lines, []string{line}...)
					}
				}
				buffer = buffer[:i]
				if len(lines) >= lineCount {
					return lines, nil
				}
			}
		}

		partialLine = string(buffer)
	}

	if partialLine != "" && len(lines) < lineCount && CheckForKeywords(partialLine, phrase) {
		lines = append(lines, []string{partialLine}...)
	}

	return lines, nil
}

// Check if one of the keywords exists in the line
func CheckForKeywords(line string, phrase string) bool {
	if len(phrase) == 0 {
		return true
	}

	lowerCaseLine := strings.ToLower(line)
	lowerCasePhrase := strings.ToLower(phrase)

	return strings.Contains(lowerCaseLine, lowerCasePhrase)
}
