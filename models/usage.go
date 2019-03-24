package models

import (
	"encoding/json"
	"net/http"
	"time"
)

// Usage type represents an action of clicking url
type Usage struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
	UrlID   string    `json:"urlId,omitempty"`
}

func (u *Usage) populate(r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(&u)
}
