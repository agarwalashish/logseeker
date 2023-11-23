package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

type BaseHandler struct {
}

// WriteJSON sends the JSON output of the http request to the client
func (bh *BaseHandler) WriteJSON(w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	// Set the headers of the http response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	b, err := json.Marshal(v) // Get the json encoding of the interface v

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"Error" : "Could not parse json interface"}`)
	} else {
		w.WriteHeader(code)
		io.WriteString(w, string(b))
	}
}
