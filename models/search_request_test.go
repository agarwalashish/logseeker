package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalSearchRequest(t *testing.T) {
	jsonString := `{"numLines":10, "keywords":"error", "filename":"log.txt"}`
	expected := SearchRequest{
		NumLines: 10,
		Keywords: "error",
		Filename: "log.txt",
	}

	var result SearchRequest
	err := json.Unmarshal([]byte(jsonString), &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
