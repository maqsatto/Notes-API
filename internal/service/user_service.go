package service

import (
	"context"
	"errors"
	"strings"

	"github.com/maqsatto/Notes-API/internal/domain"
	"github.com/maqsatto/Notes-API/internal/repository"
)

type Service struct {
	users repository.UserRepository
	// notes repository.NoteRepository
	// hasher PasswordHasher
}

func NewService(users repository.UserRepository) *Service {
	return &Service{users: users}
}

func (s *Service) Register(ctx context.Context, email, username, password string) (*domain.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	username = strings.TrimSpace(username)
	if username == "" || email == "" || password == "" {
		return nil, errors.New("Invalid input")
	}
	u := &domain.User{
		Email:    email,
		Username: username,
		Password: password,
	}
	if err := s.users.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
