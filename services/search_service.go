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

	return getLastLines(filename, lineCount)
}

// getLastLines reads the last 'lineCount' lines from a file.
func getLastLines(filename string, lineCount int) ([]string, error) {
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

	// Buffer to hold the chunk being read
	buffer := make([]byte, chunkSize)

	// Slice to store the last lines
	lines := make([]string, 0, lineCount)

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return nil, eris.New("file is empty")
	}

	// Move the file pointer to the end of the file as the newest lines are at the end of the file
	file.Seek(0, io.SeekEnd)

	// Read the file in chunks starting from the end of the file
	for fileSize > 0 {
		var offset int64
		// If the remaining file size is smaller than the chunk size,
		// adjust the chunk size and seek position accordingly.
		if fileSize < chunkSize {
			offset, err = file.Seek(-fileSize, io.SeekCurrent)
			buffer = buffer[:fileSize]
		} else {
			offset, err = file.Seek(-chunkSize, io.SeekCurrent)
		}
		if err != nil {
			return nil, err
		}

		// Read the chunk into the buffer
		_, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		// Process the chunk to extract lines
		for i := len(buffer) - 1; i >= 0; i-- {
			// Check for newline characters to identify lines
			if buffer[i] == '\n' {
				// Prepend the line to the lines slice
				lines = append([]string{string(buffer[i+1:])}, lines...)
				buffer = buffer[:i]

				// Stop if required number of lines have been collected
				if len(lines) >= lineCount {
					return lines, nil
				}
			}
		}

		// Seek back to the start of the chunk
		file.Seek(offset, io.SeekStart)

		// Update the remaining file size
		fileSize -= chunkSize
	}

	// If there's any remaining data in the buffer, treat it as the first line.
	if len(buffer) > 0 {
		lines = append([]string{string(buffer)}, lines...)
	}

	return lines, nil
}
