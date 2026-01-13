package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool

	if err := r.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}
func (r *UserRepo) CreateUser(ctx context.Context, user *domain.User) error {
	exists, err := r.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("email already exists")
	}
	query := `INSERT INTO users
    				(username, email, password, created_at, updated_at, deleted_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	now := time.Now()
	if _, err := r.db.ExecContext(ctx, query, &user.Username, &user.Email, &user.Password, now, now, nil); err != nil {
		return err
	}
	return nil
}
