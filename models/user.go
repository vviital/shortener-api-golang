package models

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// User type represents user
type User struct {
	Login      string    `json:"login"`
	Password   string    `json:"password,omitempty"`
	ID         string    `json:"id"`
	Created    time.Time `json:"created"`
	LinksCount int64     `json:"linksCount"`
	Links      []*Link   `json:"links,omitempty"`
}

// NewUserFromRequest creates new user fields from request.
func NewUserFromRequest(r *http.Request) (User, error) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	return user, err
}

func NewUserFromContext(ctx context.Context) (*User, error) {
	object := ctx.Value("user")
	err := errors.New("User is not specified in the context")

	if object == nil {
		return nil, err
	}

	user, ok := object.(*User)

	if !ok {
		return nil, err
	}

	return user, nil
}

// CleanPrivateFields removes private fields (e.g. Password) from user object
func (user *User) CleanPrivateFields() {
	user.Password = ""
}
