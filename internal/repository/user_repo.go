package repository

import (
	"context"
	"database/sql"

	"github.com/maqsatto/Notes-API/internal/domain"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
// 	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 and deleted_at IS NULL)`
// 	var exists bool
// 	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
// 	if err != nil {
// 		return false, err
// 	}
// 	return exists, nil
// }

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	q := `
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, q, user.Email, user.Username, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}
