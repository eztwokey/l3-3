package logic

import (
	"context"
	"errors"
	"strings"

	"github.com/eztwokey/l3-3/internal/models"
	"github.com/eztwokey/l3-3/internal/storage"
)

var (
	ErrBadRequest = errors.New("bad request")
)

func (l *Logic) CreateComment(ctx context.Context, req models.CreateCommentRequest) (models.Comment, error) {
	author := strings.TrimSpace(req.Author)
	text := strings.TrimSpace(req.Text)

	if author == "" || text == "" {
		return models.Comment{}, ErrBadRequest
	}

	// Если указан parent_id — проверяем, что родитель существует
	if req.ParentID != nil {
		if _, err := l.store.GetByID(ctx, *req.ParentID); err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				return models.Comment{}, ErrBadRequest
			}
			return models.Comment{}, err
		}
	}

	comment, err := l.store.Create(ctx, models.Comment{
		ParentID: req.ParentID,
		Author:   author,
		Text:     text,
	})
	if err != nil {
		l.logger.Error("create comment failed", "err", err)
		return models.Comment{}, err
	}

	l.logger.Info("comment created", "id", comment.ID, "parent_id", comment.ParentID)
	return comment, nil
}

// GetTree возвращает комментарий с id и всех его потомков в виде дерева.
func (l *Logic) GetTree(ctx context.Context, id int) (models.Comment, error) {
	if id <= 0 {
		return models.Comment{}, ErrBadRequest
	}

	flat, err := l.store.GetTree(ctx, id)
	if err != nil {
		l.logger.Error("get tree failed", "id", id, "err", err)
		return models.Comment{}, err
	}

	if len(flat) == 0 {
		return models.Comment{}, storage.ErrNotFound
	}

	// Собираем дерево из плоского списка
	tree := buildTree(flat, id)
	return tree, nil
}

func (l *Logic) ListComments(ctx context.Context, params models.ListParams) (models.CommentList, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 || params.PageSize > 100 {
		params.PageSize = 20
	}

	// Полнотекстовый поиск — отдельная ветка
	if params.Search != "" {
		comments, total, err := l.store.Search(ctx, params.Search, params.Page, params.PageSize)
		if err != nil {
			l.logger.Error("search failed", "query", params.Search, "err", err)
			return models.CommentList{}, err
		}

		return models.CommentList{
			Comments:   comments,
			Total:      total,
			Page:       params.Page,
			PageSize:   params.PageSize,
			TotalPages: totalPages(total, params.PageSize),
		}, nil
	}

	comments, total, err := l.store.List(ctx, params)
	if err != nil {
		l.logger.Error("list comments failed", "err", err)
		return models.CommentList{}, err
	}

	return models.CommentList{
		Comments:   comments,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages(total, params.PageSize),
	}, nil
}

func (l *Logic) DeleteComment(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrBadRequest
	}

	if err := l.store.Delete(ctx, id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return storage.ErrNotFound
		}
		l.logger.Error("delete comment failed", "id", id, "err", err)
		return err
	}

	l.logger.Info("comment deleted", "id", id)
	return nil
}

// buildTree собирает дерево из плоского списка комментариев.
// rootID — id корневого комментария, от которого строится дерево.
func buildTree(flat []models.Comment, rootID int) models.Comment {
	// Индексируем по ID для быстрого доступа
	byID := make(map[int]*models.Comment, len(flat))
	for i := range flat {
		flat[i].Children = []models.Comment{}
		byID[flat[i].ID] = &flat[i]
	}

	// Привязываем детей к родителям
	for i := range flat {
		if flat[i].ParentID != nil && *flat[i].ParentID != 0 {
			parent, ok := byID[*flat[i].ParentID]
			if ok && flat[i].ID != rootID {
				parent.Children = append(parent.Children, flat[i])
			}
		}
	}

	root := byID[rootID]
	if root == nil {
		return models.Comment{}
	}
	return *root
}

func totalPages(total, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	pages := total / pageSize
	if total%pageSize > 0 {
		pages++
	}
	return pages
}
