package service

import (
	"context"

	"github.com/maqsatto/Notes-API/internal/domain"
)


type UserService interface {
	Register(ctx context.Context, email, username, password string) (*domain.User, error)
	// Authenticate(ctx context.Context, email, password string) (*domain.User, error)

	// GetByID(ctx context.Context, id uint64) (*domain.User, error)
	// GetByEmail(ctx context.Context, email string) (*domain.User, error)
	// GetByUsername(ctx context.Context, username string) (*domain.User, error)

	// UpdateProfile(ctx context.Context, user *domain.User) error
	// ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error

	// SoftDelete(ctx context.Context, id uint64) error
	// HardDelete(ctx context.Context, id uint64) error

	// List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
}
