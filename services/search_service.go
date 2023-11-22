package services

import (
	"io"
	"os"
	"strings"
	"voltron/models"

	"github.com/rotisserie/eris"
)

type SearchServiceInterface interface {
	Search(request *models.SearchRequest) ([]string, error)
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

func (ss *SearchService) Search(request *models.SearchRequest) ([]string, error) {
	if request == nil || request.Filename == "" {
		return nil, eris.Errorf("missing filename")
	}

	lineCount := defaultLineCount
	if request.NumLines > 0 {
		lineCount = request.NumLines
	}

	filename := request.Filename
	filename = strings.Replace(filename, "/var/log/", "logs/", 1)

	keywords := []string{}
	if request.Keywords != "" {
		keywords = strings.Split(request.Keywords, " ")
	}

	return SearchFile(filename, lineCount, keywords)
}

// SearchFile searches for keywords in the lines from a file
func SearchFile(filename string, lineCount int, keywords []string) ([]string, error) {
	// Check if file exists before opening it
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil, eris.New("file does not exist")
	}

	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, eris.Wrap(err, "error opening file")
	}
	defer file.Close()

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return nil, eris.New("file is empty")
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
			return nil, err
		}
		_, err = file.Read(buffer)
		if err != nil {
			return nil, err
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
					if checkForKeywords(line, keywords) {
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

	if partialLine != "" && len(lines) < lineCount && checkForKeywords(partialLine, keywords) {
		lines = append(lines, []string{partialLine}...)
	}

	return lines, nil
}

// Check if one of the keywords exists in the line
func checkForKeywords(line string, keywords []string) bool {
	if len(keywords) == 0 {
		return true
	}
	m := map[string]bool{}
	for _, keyword := range keywords {
		m[strings.ToLower(keyword)] = false
	}
	words := strings.Split(line, " ")
	for _, word := range words {
		if _, ok := m[strings.ToLower(word)]; ok {
			return true
		}
	}

	return false
}
