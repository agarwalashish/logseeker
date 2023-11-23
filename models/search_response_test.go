package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalSearchResponse(t *testing.T) {
	jsonString := `{"lines":["Nov 18 06:28:45 server kernel[7616]: User logged in"]}`
	expected := SearchResponse{
		Lines: []string{
			"Nov 18 06:28:45 server kernel[7616]: User logged in",
		},
	}

	var result SearchResponse
	err := json.Unmarshal([]byte(jsonString), &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
