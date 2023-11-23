package models

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalError(t *testing.T) {
	jsonString := `{"Message":"missing filename", "Code": 400}`
	expected := Error{
		Message: "missing filename",
		Code:    http.StatusBadRequest,
	}

	var result Error
	err := json.Unmarshal([]byte(jsonString), &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
