package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	if err := r.db.QueryRowContext(ctx, query, user.Email, user.Username, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) Update(ctx context.Context, id uint64, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, username = $2
		WHERE id =$3 AND deleted_at IS NULL
		RETURNING updated_at
	`
	if err := r.db.QueryRowContext(ctx, query, user.Email, user.Username, id).Scan(&user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, id uint64, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1
		WHERE id = $2 AND deleted_at IS NULL
		RETURNING updated_at
	`
	var updated_at time.Time
	if err := r.db.QueryRowContext(ctx, query, hashedPassword, id).Scan(&updated_at); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *UserRepo) SoftDelete(ctx context.Context, id uint64) error {
	query := `
		UPDATE users
		SET deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING deleted_at
	`
	var deleted_at time.Time
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&deleted_at); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *UserRepo) HardDelete(ctx context.Context, id uint64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uint64) (*domain.User, error) {
	query := `
		SELECT id, email, username, password, created_at, updated_at FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	var getUser domain.User

	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&getUser.ID, &getUser.Email, &getUser.Username, &getUser.Password, &getUser.CreatedAt, &getUser.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &getUser, nil
}

func (r *UserRepo) GetByIDAny(ctx context.Context, id uint64) (*domain.User, error) {
	query := `
		SELECT id, email, username, password, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1
	`
	var u domain.User

	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.Username, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}



func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password, created_at, updated_at FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`
	var getUser domain.User

	if err := r.db.QueryRowContext(ctx, query, email).Scan(
		&getUser.ID, &getUser.Email, &getUser.Username, &getUser.Password, &getUser.CreatedAt, &getUser.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &getUser, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password, created_at, updated_at FROM users
		WHERE username = $1 AND deleted_at IS NULL
	`
	var getUser domain.User

	if err := r.db.QueryRowContext(ctx, query, username).Scan(
		&getUser.ID, &getUser.Email, &getUser.Username, &getUser.Password, &getUser.CreatedAt, &getUser.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &getUser, nil
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 and deleted_at IS NULL)`
	var exists bool
	if err := r.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 and deleted_at IS NULL)`
	var exists bool
	if err := r.db.QueryRowContext(ctx, query, username).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, email, username, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0, limit)
	for rows.Next() {
		user := new(domain.User)
		if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepo) Count(ctx context.Context) (uint64, error) {
	query := `
		SELECT COUNT(*) FROM users;
	`
	var count uint64
	if err := r.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
