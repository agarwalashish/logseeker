package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchFile(t *testing.T) {
	filename := "../tests/data/log_file.log"

	tests := []struct {
		expectedLines int
		numLines      int
		keywords      string
	}{
		{
			expectedLines: 10,
			numLines:      10,
			keywords:      "user",
		},
		{
			expectedLines: 90,
			numLines:      500,
			keywords:      "Error encountered",
		},
		{
			expectedLines: 10,
			numLines:      10,
		},
	}

	for _, ts := range tests {
		lines, err := SearchFile(filename, ts.numLines, ts.keywords)
		assert.Nil(t, err)
		assert.Equal(t, ts.expectedLines, len(lines))
	}
}

func TestSearchFileDoesNotExist(t *testing.T) {
	filename := "/var/log/invalid_file.txt"
	lines, err := SearchFile(filename, 10, "dummy phrase")
	assert.Nil(t, lines)
	assert.Equal(t, "file does not exist", err.Message)
}

func TestSearchEmptyFile(t *testing.T) {
	filename := "../tests/data/empty_file.log"
	lines, err := SearchFile(filename, 10, "dummy phrase")
	assert.Nil(t, lines)
	assert.Equal(t, "file is empty", err.Message)
}

func TestCheckForKeywords(t *testing.T) {
	tests := []struct {
		line     string
		phrase   string
		expected bool
	}{
		{"This is a test line", "test", true},
		{"Another line", "nope", false},
		{"Case Insensitive TEST", "test", true},
		{"Empty phrase should return true", "", true},
		{"", "no match", false},
		{"Exact match", "Exact match", true},
	}

	for _, test := range tests {
		result := CheckForKeywords(test.line, test.phrase)
		assert.Equal(t, test.expected, result)
	}
}
