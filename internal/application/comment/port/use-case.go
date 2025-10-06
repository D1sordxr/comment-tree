package port

import (
	"context"
	"github.com/D1sordxr/comment-tree/internal/application/comment/input"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/model"
)

type UseCase interface {
	Create(ctx context.Context, i input.CreateInput) (*model.Comment, error)
	GetTreeByDestination(ctx context.Context, dest string) (model.CommentTree, error)

	SearchSimilar(ctx context.Context, i input.SearchSimilarInput) (model.RawComments, error)
}
