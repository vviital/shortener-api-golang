package models

import (
	"encoding/json"
	"net/http"
)

// Link struct represents link
type Link struct {
	Created     string  `json:"created"`
	ID          string  `json:"id"`
	URL         string  `json:"url"`
	UsagesCount int64   `json:"usagesCount"`
	UserID      string  `json:"userId"`
	Usages      []Usage `json:"usages,omitempty"`
}

func (l *Link) Populate(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(&l)
}
