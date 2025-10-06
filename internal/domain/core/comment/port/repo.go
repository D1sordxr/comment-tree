package port

import (
	"context"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/model"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/params"
)

type Repository interface {
	Create(ctx context.Context, p params.Create) (*model.Comment, error)
	GetByDestination(ctx context.Context, dest string) ([]model.Comment, error)
	GetWithPagination(ctx context.Context, p params.GetWithPagination) ([]model.Comment, error)
	Delete(ctx context.Context, id int) error
}
