package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/eztwokey/l3-3/internal/models"
)

func (s *Storage) Create(ctx context.Context, c models.Comment) (models.Comment, error) {
	var id int
	var createdAt sql.NullTime

	err := s.db.QueryRowContext(ctx,
		`INSERT INTO comments (parent_id, author, text)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		c.ParentID, c.Author, c.Text,
	).Scan(&id, &createdAt)

	if err != nil {
		return models.Comment{}, fmt.Errorf("insert comment: %w", err)
	}

	c.ID = id
	if createdAt.Valid {
		c.CreatedAt = createdAt.Time
	}

	return c, nil
}

func (s *Storage) GetByID(ctx context.Context, id int) (models.Comment, error) {
	var c models.Comment
	var parentID sql.NullInt64

	err := s.db.QueryRowContext(ctx,
		`SELECT id, parent_id, author, text, created_at
		 FROM comments WHERE id = $1`,
		id,
	).Scan(&c.ID, &parentID, &c.Author, &c.Text, &c.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Comment{}, ErrNotFound
		}
		return models.Comment{}, fmt.Errorf("get comment: %w", err)
	}

	if parentID.Valid {
		pid := int(parentID.Int64)
		c.ParentID = &pid
	}

	return c, nil
}

func (s *Storage) GetTree(ctx context.Context, id int) ([]models.Comment, error) {
	rows, err := s.db.QueryContext(ctx,
		`WITH RECURSIVE tree AS (
			SELECT id, parent_id, author, text, created_at
			FROM comments WHERE id = $1
			UNION ALL
			SELECT c.id, c.parent_id, c.author, c.text, c.created_at
			FROM comments c JOIN tree t ON c.parent_id = t.id
		)
		SELECT id, parent_id, author, text, created_at FROM tree
		ORDER BY created_at`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("get tree: %w", err)
	}
	defer func() { _ = rows.Close() }()

	return scanComments(rows)
}

func (s *Storage) List(ctx context.Context, params models.ListParams) ([]models.Comment, int, error) {
	// Определяем условие фильтрации
	where := "WHERE parent_id IS NULL"
	args := []interface{}{}
	argIdx := 1

	if params.ParentID != nil {
		where = fmt.Sprintf("WHERE parent_id = $%d", argIdx)
		args = append(args, *params.ParentID)
		argIdx++
	}

	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM comments %s", where)
	err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count comments: %w", err)
	}

	orderBy := "created_at DESC"
	if params.SortBy == "created_at" && params.SortDir == "asc" {
		orderBy = "created_at ASC"
	}

	offset := (params.Page - 1) * params.PageSize

	query := fmt.Sprintf(
		"SELECT id, parent_id, author, text, created_at FROM comments %s ORDER BY %s LIMIT $%d OFFSET $%d",
		where, orderBy, argIdx, argIdx+1,
	)
	args = append(args, params.PageSize, offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list comments: %w", err)
	}
	defer func() { _ = rows.Close() }()

	comments, err := scanComments(rows)
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (s *Storage) Search(ctx context.Context, query string, page, pageSize int) ([]models.Comment, int, error) {
	var total int
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM comments
		 WHERE search_vector @@ plainto_tsquery('russian', $1)`,
		query,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("search count: %w", err)
	}

	offset := (page - 1) * pageSize

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, parent_id, author, text, created_at FROM comments
		 WHERE search_vector @@ plainto_tsquery('russian', $1)
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		query, pageSize, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("search comments: %w", err)
	}
	defer func() { _ = rows.Close() }()

	comments, err := scanComments(rows)
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx,
		`DELETE FROM comments WHERE id = $1`, id,
	)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func scanComments(rows *sql.Rows) ([]models.Comment, error) {
	var comments []models.Comment

	for rows.Next() {
		var c models.Comment
		var parentID sql.NullInt64

		if err := rows.Scan(&c.ID, &parentID, &c.Author, &c.Text, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}

		if parentID.Valid {
			pid := int(parentID.Int64)
			c.ParentID = &pid
		}

		comments = append(comments, c)
	}

	return comments, rows.Err()
}
