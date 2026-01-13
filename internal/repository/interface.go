package repository

import (
	"context"
	"errors"

	"github.com/maqsatto/Notes-API/internal/domain"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, note *domain.Note) error
	UpdateNote(ctx context.Context, note *domain.Note) error
	DeleteNote(ctx context.Context, id uint64) error
	GetNoteByID(ctx context.Context, id uint64) (*domain.Note, error)
	GetNoteByTitle(ctx context.Context, title string) (*domain.Note, error)
	GetNoteByTags(ctx context.Context, tags []string) (*domain.Note, error)
	GetNotesByUserID(ctx context.Context, userID uint64) ([]*domain.Note, error)
	GetAllNotes(ctx context.Context) ([]*domain.Note, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint64) error
	GetUserByID(ctx context.Context, id uint64) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
