package models

type SearchRequest struct {
	NumLines int    `json:"numLines"`
	Keywords string `json:"keywords"`
	Filename string `json:"filename"`
}
