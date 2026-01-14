package service

import (
	"context"

	"github.com/maqsatto/Notes-API/internal/domain"
)

type noteService interface {
	Create(ctx context.Context, userID uint64, title, content string, tags []string) (*domain.Note, error)
	Update(ctx context.Context, userID, noteID uint64, title, content string, tags []string) (*domain.Note, error)
	Delete(ctx context.Context, userID, noteID uint64) error
	PermanentDelete(ctx context.Context, userID, noteID uint64) error

	GetByID(ctx context.Context, userID, noteID uint64) (*domain.Note, error)
	ListByUser(ctx context.Context, userID uint64, limit, offset int) ([]*domain.Note, int64, error)
	ListByTags(ctx context.Context, userID uint64, tags []string, limit, offset int) ([]*domain.Note, int64, error)
	SearchByTitle(ctx context.Context, userID uint64, query string, limit, offset int) ([]*domain.Note, int64, error)
	SearchByContent(ctx context.Context, userID uint64, query string, limit, offset int) ([]*domain.Note, int64, error)

	GetUserNoteCount(ctx context.Context, userID uint64) (int64, error)
}

type userService interface {
	Register(ctx context.Context, username, email, password string) (*domain.User, error)
	Login(ctx context.Context, emailOrUsername, password string) (*domain.User, string, error)
	UpdateProfile(ctx context.Context, userID uint64, username, email string) (*domain.User, error)
	ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error

	DeleteAccount(ctx context.Context, userID uint64) error
	PermanentDeleteAccount(ctx context.Context, userID uint64) error

	GetByID(ctx context.Context, userID uint64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)

	ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
	IsUsernameTaken(ctx context.Context, username string) (bool, error)

	GetTotalUserCount(ctx context.Context) (uint64, error)
}
