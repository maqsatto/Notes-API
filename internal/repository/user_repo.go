package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/maqsatto/Notes-API/internal/domain"
)

type userRepo struct {
	ctx context.Context
	db  *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{
		ctx: context.Background(),
		db:  db,
	}
}

func (r *userRepo) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (username, email, password, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6)`
	timeDateNow := time.Now()
	r.db.QueryRow(query, &user.Username, &user.Username, &user.Password, timeDateNow, timeDateNow, "NULL")
	return nil
}
