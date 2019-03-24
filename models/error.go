package models

// Error struct
type Error struct {
	Message string `json:"message"`
}

// NewError creates new Error object
func NewError(message string) Error {
	return Error{message}
}
