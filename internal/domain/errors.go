package domain

import "errors"

// Common / Base errors

var (
	ErrInvalidID      = errors.New("invalid id")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrConflict       = errors.New("resource conflict")
	ErrInternal       = errors.New("internal server error")
	ErrNotImplemented = errors.New("not implemented")
)

// User errors

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")

	ErrInvalidUser     = errors.New("invalid user data")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")

	// Specific password validation errors
	ErrPasswordTooShort         = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong          = errors.New("password must not exceed 128 characters")
	ErrPasswordMissingUppercase = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingLowercase = errors.New("password must contain at least one lowercase letter")
	ErrPasswordMissingDigit     = errors.New("password must contain at least one digit")
	ErrPasswordMissingSpecial   = errors.New("password must contain at least one special character")

	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordMismatch   = errors.New("password mismatch")

	ErrUserDeleted = errors.New("user is deleted")
)

// Note errors

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrInvalidNote  = errors.New("invalid note data")

	ErrNoteTitleEmpty   = errors.New("note title is empty")
	ErrNoteContentEmpty = errors.New("note content is empty")
	ErrNoteTooLarge     = errors.New("note content too large")

	ErrNoteAccessDenied = errors.New("access to note denied")
	ErrNoteDeleted      = errors.New("note is deleted")

	ErrInvalidTags = errors.New("invalid tags")
	ErrTooManyTags = errors.New("too many tags")
)

// Repository / persistence errors

var (
	ErrAlreadyDeleted  = errors.New("already deleted")
	ErrNothingToUpdate = errors.New("nothing to update")

	ErrDatabase    = errors.New("database error")
	ErrTransaction = errors.New("transaction error")
)

// Auth / JWT errors

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrExpiredToken   = errors.New("token expired")
	ErrMissingToken   = errors.New("missing token")
	ErrMalformedToken = errors.New("malformed token")

	ErrInvalidIssuer        = errors.New("invalid token issuer")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)

// Pagination / filtering

var (
	ErrInvalidLimit       = errors.New("invalid limit")
	ErrInvalidOffset      = errors.New("invalid offset")
	ErrInvalidSearchQuery = errors.New("invalid search query")
)

// Service-level errors

var (
	ErrOperationNotAllowed = errors.New("operation not allowed")
	ErrStateViolation      = errors.New("state violation")
)
