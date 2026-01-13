package repository

import (
	"context"

	"github.com/maqsatto/Notes-API/internal/domain"
)

// Common repository errors
// var (
// 	ErrNotFound          = errors.New("record not found")
// 	ErrDuplicateEmail    = errors.New("email already exists")
// 	ErrDuplicateUsername = errors.New("username already exists")
// 	ErrInvalidData       = errors.New("invalid data")
// )

// type Transactor interface {
// 	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
// }

type NoteRepository interface {
	Create(ctx context.Context, note *domain.Note) error
	Update(ctx context.Context, note *domain.Note) error
	SoftDelete(ctx context.Context, id uint64) error
	HardDelete(ctx context.Context, id uint64) error

	GetByID(ctx context.Context, id uint64) (*domain.Note, error)
	ListByUserID(ctx context.Context, userID uint64, limit, offset int) ([]*domain.Note, int64, error)
	ListByTags(ctx context.Context, userID uint64, tags []string, limit, offset int) ([]*domain.Note, int64, error)
	SearchByTitle(ctx context.Context, userID uint64, titleQuery string, limit, offset int) ([]*domain.Note, int64, error)
	SearchByContent(ctx context.Context, userID uint64, contentQuery string, limit, offset int) ([]*domain.Note, int64, error)

	CountByUserID(ctx context.Context, userID uint64) (int64, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, id uint64, user *domain.User) error
	UpdatePassword(ctx context.Context, id uint64, hashedPassword string) error

	SoftDelete(ctx context.Context, id uint64) error

	HardDelete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.User, error)

	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)

	// List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	Count(ctx context.Context) (uint64, error)
}

// type TagRepository interface {
// 	GetOrCreateTags(ctx context.Context, tagNames []string) ([]int64, error)
// 	ListPopularTags(ctx context.Context, limit int) ([]string, error)
// 	ListUserTags(ctx context.Context, userID uint64) ([]string, error)
// 	DeleteUnusedTags(ctx context.Context) error
// }
