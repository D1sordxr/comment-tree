package repo

import (
	"context"
	"fmt"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/model"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/params"
	"github.com/D1sordxr/comment-tree/internal/infrastructure/storage/postgres/repositories/comment/converters"
	"github.com/D1sordxr/comment-tree/internal/infrastructure/storage/postgres/repositories/comment/gen"
	"github.com/wb-go/wbf/dbpg"
)

type Repository struct {
	executor *dbpg.DB
	queries  *gen.Queries
}

func New(executor *dbpg.DB) *Repository {
	return &Repository{
		executor: executor,
		queries:  gen.New(executor.Master),
	}

}

func (r *Repository) Create(ctx context.Context, p params.Create) (*model.Comment, error) {
	const op = "comment.Repository.Create"

	rawComment, err := r.queries.CreateComment(ctx, converters.ConvertCreateParams(p))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	comment := converters.ConvertGenToDomain(rawComment)
	return &comment, nil
}

func (r *Repository) GetByDestination(ctx context.Context, dest string) ([]model.Comment, error) {
	const op = "comment.Repository.GetByDestination"

	rawComments, err := r.queries.GetCommentsByDestination(ctx, dest)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converters.ConvertGenSliceToDomain(rawComments), nil
}

func (r *Repository) GetWithPagination(ctx context.Context, p params.GetWithPagination) ([]model.Comment, error) {
	const op = "comment.Repository.GetBySource"

	tx, err := r.executor.Master.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	qtx := r.queries.WithTx(tx)
	ids, err := qtx.GetRootCommentIDsWithPagination(ctx, converters.ConvertGetRootIDsWithPaginationParams(p))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rawComments, err := qtx.GetCommentsWithChildren(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	comments := make(model.RawComments, len(rawComments))
	for i, rawComment := range rawComments {
		comments[i] = converters.ConvertGetCommentsWithChildrenRowToModel(rawComment)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	const op = "comment.Repository.Delete"

	if err := r.queries.DeleteCommentByID(ctx, int32(id)); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
