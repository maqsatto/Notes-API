package domain

import "errors"

// User related errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUser       = errors.New("invalid user data")
)

// Note related errors
var (
	ErrNoteNotFound = errors.New("note not found")
	ErrInvalidNote  = errors.New("invalid note data")
	ErrUnauthorized = errors.New("unauthorized access to note")
)

// Common errors
var (
	ErrInvalidID = errors.New("invalid id")
	ErrInternal  = errors.New("internal server error")
)
