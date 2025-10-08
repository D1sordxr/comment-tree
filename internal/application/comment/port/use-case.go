package port

import (
	"context"
	"github.com/D1sordxr/comment-tree/internal/application/comment/input"
	"github.com/D1sordxr/comment-tree/internal/application/comment/output"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/model"
)

type UseCase interface {
	Create(ctx context.Context, input input.CreateInput) (*output.CreateOutput, error)
	GetByIDWithChildren(ctx context.Context, id int) (*model.Comment, error)
	GetTreeByDestination(ctx context.Context, dest string) (*output.GetTreeOutput, error)
	GetTreeWithPagination(ctx context.Context, i input.GetTreeWithPagination) (*output.GetTreeWithPaginationOutput, error)
	GetCommentsWithPagination(ctx context.Context, i input.GetCommentsWithPagination) (*output.GetCommentsWithPaginationOutput, error)
	Delete(ctx context.Context, id int) (*output.DeleteOutput, error)
	SearchSimilar(ctx context.Context, i input.SearchSimilarInput) (*output.SearchSimilarOutput, error)
}
