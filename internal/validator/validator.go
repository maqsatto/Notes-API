package validator

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

var (
	ErrInvalidEmailAddress = errors.New("invalid email address")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrEmptyValue          = errors.New("value is empty")
	ErrInvalidUsername     = errors.New("invalid username")
	ErrInvalidNoteTitle    = errors.New("invalid note title")
	ErrInvalidNoteContent  = errors.New("invalid note content")
	ErrInvalidID           = errors.New("invalid id")
	ErrValueTooShort       = errors.New("value is too short")
	ErrValueTooLong        = errors.New("value is too long")
)

const (
	MinPasswordLength    = 8
	MaxPasswordLength    = 128
	MinUsernameLength    = 3
	MaxUsernameLength    = 50
	MaxNoteTitleLength   = 200
	MaxNoteContentLength = 50000
)

func ValidateUserRegister(email, username, password string) error {
	if _, err := IsValidEmail(email); err != nil {
		return err
	}
	if _, err := IsValidPassword(password); err != nil {
		return err
	}

	if username != "" {
		if _, err := IsValidUsername(username); err != nil {
			return err
		}
	}
	return nil
}

func ValidateUserLogin(email, password string) error {
	if _, err := IsValidEmail(email); err != nil {
		return err
	}

	if _, err := IsEmptyString(password); err != nil {
		return ErrEmptyValue
	}

	return nil
}

func IsValidNote(title, content string) error {
	if _, err := IsValidTitle(title); err != nil {
		return err
	}

	if _, err := IsValidContent(content); err != nil {
		return err
	}

	return nil
}
func IsValidEmail(email string) (bool, error) {

	if len(email) > 255 {
		return false, ErrInvalidEmailAddress
	}

	if b, err := IsEmptyString(email); err != nil {
		return b, ErrInvalidEmailAddress
	}

	email = strings.TrimSpace(email)
	addr, err := mail.ParseAddress(email)

	if err != nil || addr == nil || addr.Address != email {
		return false, ErrInvalidEmailAddress
	}

	parts := strings.Split(addr.Address, "@")

	if len(parts) != 2 {
		return false, ErrInvalidEmailAddress
	}

	user, domain := parts[0], parts[1]

	if !strings.Contains(domain, ".") {
		return false, ErrInvalidEmailAddress
	}

	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false, ErrInvalidEmailAddress
	}

	if len(user) < 1 || len(domain) < 3 {
		return false, ErrInvalidEmailAddress
	}

	return true, nil
}

func IsEmptyString(value string) (bool, error) {
	if value == "" || len(strings.TrimSpace(value)) == 0 {
		return false, ErrEmptyValue
	}
	return true, nil
}

func IsValidPassword(password string) (bool, error) {

	if _, err := IsEmptyString(password); err != nil {
		return false, ErrInvalidPassword
	}

	if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
		return false, ErrInvalidPassword
	}

	var (
		hasUppercase,
		hasLowercase,
		hasDigit,
		hasSpecialChar bool
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}
	if !hasUppercase || !hasLowercase || !hasDigit || !hasSpecialChar {
		return false, ErrInvalidPassword
	}

	return true, nil
}

func IsValidUsername(username string) (bool, error) {
	if _, err := IsEmptyString(username); err != nil {
		return false, ErrInvalidUsername
	}
	if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
		return false, ErrInvalidUsername
	}
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validUsername.MatchString(username) {
		return false, ErrInvalidUsername
	}
	return true, nil
}

func IsValidTitle(title string) (bool, error) {
	if _, err := IsEmptyString(title); err != nil {
		return false, ErrInvalidNoteTitle
	}
	if len(title) > MaxNoteTitleLength {
		return false, ErrInvalidNoteTitle
	}
	return true, nil
}

func IsValidContent(content string) (bool, error) {
	if _, err := IsEmptyString(content); err != nil {
		return false, ErrInvalidNoteTitle
	}
	if len(content) > MaxNoteContentLength {
		return false, ErrInvalidNoteTitle
	}
	return true, nil
}
