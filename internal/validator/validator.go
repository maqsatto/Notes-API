package validator

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode"

	"github.com/maqsatto/Notes-API/internal/domain"
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
		return domain.ErrInvalidInput
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
		return false, domain.ErrInvalidEmail
	}

	if b, err := IsEmptyString(email); err != nil {
		return b, domain.ErrInvalidEmail
	}

	email = strings.TrimSpace(email)
	addr, err := mail.ParseAddress(email)

	if err != nil || addr == nil || addr.Address != email {
		return false, domain.ErrInvalidEmail
	}

	parts := strings.Split(addr.Address, "@")

	if len(parts) != 2 {
		return false, domain.ErrInvalidEmail
	}

	user, domainPart := parts[0], parts[1]

	if !strings.Contains(domainPart, ".") {
		return false, domain.ErrInvalidEmail
	}

	if strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") {
		return false, domain.ErrInvalidEmail
	}

	if len(user) < 1 || len(domainPart) < 3 {
		return false, domain.ErrInvalidEmail
	}

	return true, nil
}

func IsEmptyString(value string) (bool, error) {
	if value == "" || len(strings.TrimSpace(value)) == 0 {
		return false, domain.ErrInvalidInput
	}
	return true, nil
}

func IsValidPassword(password string) (bool, error) {
	if _, err := IsEmptyString(password); err != nil {
		return false, domain.ErrInvalidPassword
	}

	if len(password) < MinPasswordLength {
		return false, domain.ErrPasswordTooShort
	}

	if len(password) > MaxPasswordLength {
		return false, domain.ErrPasswordTooLong
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

	if !hasUppercase {
		return false, domain.ErrPasswordMissingUppercase
	}
	if !hasLowercase {
		return false, domain.ErrPasswordMissingLowercase
	}
	if !hasDigit {
		return false, domain.ErrPasswordMissingDigit
	}
	if !hasSpecialChar {
		return false, domain.ErrPasswordMissingSpecial
	}

	return true, nil
}

func IsValidUsername(username string) (bool, error) {
	if _, err := IsEmptyString(username); err != nil {
		return false, domain.ErrInvalidUsername
	}
	if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
		return false, domain.ErrInvalidUsername
	}
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validUsername.MatchString(username) {
		return false, domain.ErrInvalidUsername
	}
	return true, nil
}

func IsValidTitle(title string) (bool, error) {
	if _, err := IsEmptyString(title); err != nil {
		return false, domain.ErrNoteTitleEmpty
	}
	if len(title) > MaxNoteTitleLength {
		return false, domain.ErrInvalidNote
	}
	return true, nil
}

func IsValidContent(content string) (bool, error) {
	if _, err := IsEmptyString(content); err != nil {
		return false, domain.ErrNoteContentEmpty
	}
	if len(content) > MaxNoteContentLength {
		return false, domain.ErrNoteTooLarge
	}
	return true, nil
}

func ValidateUserUpdate(email, username string) error {
	if _, err := IsValidEmail(email); err != nil {
		return err
	}
	if username != "" {
		if _, err := IsValidUsername(username); err != nil {
			return err
		}
	}
	return nil
}
