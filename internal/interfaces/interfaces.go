package interfaces

import (
	"context"

	"github.com/eztwokey/l3-3/internal/models"
)

type CommentStorage interface {
	Create(ctx context.Context, comment models.Comment) (models.Comment, error)
	GetByID(ctx context.Context, id int) (models.Comment, error)
	GetTree(ctx context.Context, id int) ([]models.Comment, error)
	List(ctx context.Context, params models.ListParams) ([]models.Comment, int, error)
	Search(ctx context.Context, query string, page, pageSize int) ([]models.Comment, int, error)
	Delete(ctx context.Context, id int) error
}
