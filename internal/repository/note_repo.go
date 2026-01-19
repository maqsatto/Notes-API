package repository

import (
	"context"
	"database/sql"

	"github.com/maqsatto/Notes-API/internal/domain"
)

type NoteRepo struct {
	db *sql.DB
}

func NewNoteRepo(db *sql.DB) *NoteRepo {
	return &NoteRepo{
		db: db,
	}
}

func (r *NoteRepo) Create(ctx context.Context, note *domain.Note) error {

	query := `INSERT INTO notes (title, content, user_id, tags)
			  VALUES ($1, $2, $3, $4)
			  RETURNING id, created_at, updated_at`

	if err := r.db.QueryRowContext(ctx, query, note.Title, note.Content, note.UserID, note.Tags).
		Scan(&note.ID, &note.CreatedAt, &note.UpdatedAt); err != nil {
		return err
	}
	return nil
}



func (r *NoteRepo) Update(ctx context.Context, note *domain.Note) error {

	query := `UPDATE notes SET title = $1, content = $2, tags = $3
             WHERE id = $4 AND deleted_at IS NULL RETURNING updated_at`

	if err := r.db.QueryRowContext(ctx, query, note.Title, note.Content, note.Tags, note.ID).
		Scan(&note.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (r *NoteRepo) SoftDelete(ctx context.Context, id uint64) error {

	query := `UPDATE notes SET deleted_at = now()
             WHERE id = $1 AND deleted_at IS NULL RETURNING deleted_at`

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (r *NoteRepo) HardDelete(ctx context.Context, id uint64) error {
	query := `DELETE FROM notes WHERE id = $1`
	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}

func (r *NoteRepo) GetByID(ctx context.Context, id uint64) (*domain.Note, error) {

	query := `SELECT id, user_id, title, content, created_at, updated_at, tags FROM notes
                WHERE id = $1 AND deleted_at IS NULL`

	var note domain.Note
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.Tags); err != nil {
		return nil, err
	}

	return &note, nil
}
func (r *NoteRepo) ListByUserID(
	ctx context.Context,
	userID uint64,
	limit, offset int,
) ([]*domain.Note, int64, error) {

	countQuery := `
		SELECT COUNT(*)
		FROM notes
		WHERE user_id = $1
		  AND deleted_at IS NULL
	`

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT id, user_id, title, content, created_at, updated_at, tags
		FROM notes
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, dataQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notes := make([]*domain.Note, 0, limit)

	for rows.Next() {
		var note domain.Note
		if err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
			&note.Tags,
		); err != nil {
			return nil, 0, err
		}

		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

func (r *NoteRepo) ListByTags(ctx context.Context, userID uint64, tags []string, limit, offset int) ([]*domain.Note, int64, error) {
	dataQuery := `
		SELECT id, user_id, title, content, created_at, updated_at, tags
		FROM notes
		WHERE user_id = $1
		AND deleted_at IS NULL
		AND tags @> $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	countQuery := `	SELECT COUNT(*)
					FROM notes
					WHERE user_id = $1
					AND deleted_at IS NULL
					AND tags @> $2`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, userID, tags).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx, dataQuery, userID, tags, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notes := make([]*domain.Note, 0, limit)

	for rows.Next() {
		var note domain.Note
		if err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
			&note.Tags); err != nil {
			return nil, 0, err
		}
		notes = append(notes, &note)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

func (r *NoteRepo) SearchByTitle(
	ctx context.Context,
	userID uint64,
	titleQuery string,
	limit, offset int,
) ([]*domain.Note, int64, error) {

	search := "%" + titleQuery + "%"

	countQuery := `
		SELECT COUNT(*)
		FROM notes
		WHERE user_id = $1
		  AND deleted_at IS NULL
		  AND title ILIKE $2
	`

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, userID, search).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQuery := `
		SELECT id, user_id, title, content, created_at, updated_at, tags
		FROM notes
		WHERE user_id = $1
		  AND deleted_at IS NULL
		  AND title ILIKE $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, dataQuery, userID, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notes := make([]*domain.Note, 0, limit)

	for rows.Next() {
		note := new(domain.Note)
		if err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
			&note.Tags,
		); err != nil {
			return nil, 0, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

func (r *NoteRepo) SearchByContent(ctx context.Context, userID uint64, contentQuery string, limit, offset int) ([]*domain.Note, int64, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM notes
		WHERE user_id = $1
		AND deleted_at IS NULL
		AND content ILIKE $2
	`
	search := "%" + contentQuery + "%"
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, userID, search).Scan(&total); err != nil {
		return nil, 0, err
	}
	dataQuery := `
		SELECT id, user_id, title, content, created_at, updated_at, tags
		FROM notes
		WHERE user_id = $1
		  AND deleted_at IS NULL
		  AND content ILIKE $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.QueryContext(ctx, dataQuery, userID, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notes := make([]*domain.Note, 0, limit)
	for rows.Next() {
		var note domain.Note
		if err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
			&note.Tags,
		); err != nil {
			return nil, 0, err
		}

		notes = append(notes, &note)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

func (r *NoteRepo) CountByUserID(ctx context.Context, userID uint64) (int64, error) {
	query := `SELECT COUNT(*) FROM notes WHERE user_id = $1 AND deleted_at IS NULL`
	var count int64
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
