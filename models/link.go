package models

import (
	"encoding/json"
	"net/http"
	"time"
)

// Link struct represents link
type Link struct {
	Created     time.Time `json:"created"`
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	UsagesCount int64     `json:"usagesCount"`
	UserID      string    `json:"userId"`
	Usages      []Usage   `json:"usages,omitempty"`
}

func (l *Link) Populate(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(&l)
}
