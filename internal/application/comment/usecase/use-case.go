package usecase

import (
	"context"
	"fmt"
	"github.com/D1sordxr/comment-tree/internal/application/comment/input"
	"github.com/D1sordxr/comment-tree/internal/application/comment/output"
	appPorts "github.com/D1sordxr/comment-tree/internal/domain/app/ports"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/errorx"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/model"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/params"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/port"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/service"
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/vo"
	"github.com/D1sordxr/comment-tree/pkg/logger"
	"strings"
)

type UseCase struct {
	log  appPorts.Logger
	repo port.Repository
}

func New(
	log appPorts.Logger,
	repo port.Repository,
) *UseCase {
	return &UseCase{
		log:  log,
		repo: repo,
	}
}

func (uc *UseCase) Create(
	ctx context.Context,
	input input.CreateInput,
) (*output.CreateOutput, error) {
	const op = "comment.UseCase.Create"
	logFields := logger.WithFields("operation", op)

	uc.log.Info("Attempting to create a new comment", logFields()...)

	comment, err := uc.repo.Create(ctx, params.Create{
		CommentDestination: input.CommentDestination,
		Content:            input.Content,
		Author:             input.Author,
		ParentID:           input.ParentID,
	})
	if err != nil {
		uc.log.Error("Failed to create a new comment", logFields("error", err)...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	uc.log.Info("Successfully created a new comment", logFields()...)

	return &output.CreateOutput{Comment: comment}, nil
}

func (uc *UseCase) GetTreeByDestination(
	ctx context.Context,
	dest string,
) (*output.GetTreeOutput, error) {
	const op = "comment.UseCase.GetByDestination"
	logFields := logger.WithFields("operation", op)

	uc.log.Info("Attempting to get comments by destination", logFields("destination", dest)...)

	comments, err := uc.repo.GetByDestination(ctx, dest)
	if err != nil {
		uc.log.Error("Failed to get comments by destination", logFields("error", err)...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tree := service.BuildCommentTree(comments)
	uc.log.Info("Successfully built comment tree", logFields(
		"total_comments", len(comments), "tree_size", len(tree),
	)...)

	return &output.GetTreeOutput{Tree: tree}, nil
}

func (uc *UseCase) GetTreeWithPagination(
	ctx context.Context,
	i input.GetTreeWithPagination,
) (*output.GetTreeWithPaginationOutput, error) {
	const op = "comment.UseCase.GetTreeWithPagination"
	logFields := logger.WithFields("operation", op, "destination", i.CommentDestination)

	uc.log.Info("Attempting to get paginated comment tree", logFields()...)

	comments, err := uc.repo.GetWithPagination(ctx, params.GetWithPagination{
		Destination: i.CommentDestination,
		CursorID:    i.CursorID,
		Limit:       vo.DefaultLimit.Int(),
	})
	if err != nil {
		uc.log.Error("Failed to get paginated comments", logFields("error", err)...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tree := service.BuildCommentTree(comments)

	// Вычисляем next cursor
	nextCursor := 0
	if len(comments) > 0 {
		// Находим максимальный ID среди комментариев
		for _, comment := range comments {
			if comment.ID > nextCursor {
				nextCursor = comment.ID
			}
		}
	}

	uc.log.Info("Successfully built paginated comment tree", logFields(
		"total_comments", len(comments), "tree_size", len(tree), "next_cursor", nextCursor,
	)...)

	return &output.GetTreeWithPaginationOutput{
		Tree:       tree,
		NextCursor: nextCursor,
	}, nil
}

func (uc *UseCase) GetCommentsWithPagination(
	ctx context.Context,
	i input.GetCommentsWithPagination,
) (*output.GetCommentsWithPaginationOutput, error) {
	const op = "comment.UseCase.GetCommentsWithPagination"
	logFields := logger.WithFields("operation", op, "destination", i.CommentDestination, "cursor", i.CursorID)

	uc.log.Info("Attempting to get paginated comments", logFields()...)

	comments, err := uc.repo.GetWithPagination(ctx, params.GetWithPagination{
		Destination: i.CommentDestination,
		CursorID:    i.CursorID,
		Limit:       vo.DefaultLimit.Int(),
	})
	if err != nil {
		uc.log.Error("Failed to get paginated comments", logFields("error", err)...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Вычисляем next cursor
	nextCursor := 0
	if len(comments) > 0 {
		nextCursor = comments[len(comments)-1].ID
	}

	uc.log.Info("Successfully retrieved paginated comments", logFields(
		"count", len(comments), "next_cursor", nextCursor,
	)...)

	return &output.GetCommentsWithPaginationOutput{
		Comments:   comments,
		NextCursor: nextCursor,
	}, nil
}

func (uc *UseCase) Delete(
	ctx context.Context,
	id int,
) (*output.DeleteOutput, error) {
	const op = "comment.UseCase.Delete"
	logFields := logger.WithFields("operation", op, "comment_id", id)

	uc.log.Info("Attempting to delete comment", logFields()...)

	err := uc.repo.Delete(ctx, id)
	if err != nil {
		uc.log.Error("Failed to delete comment", logFields("error", err)...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	uc.log.Info("Successfully deleted comment", logFields()...)
	return &output.DeleteOutput{Success: true}, nil
}

func (uc *UseCase) SearchSimilar(
	ctx context.Context,
	i input.SearchSimilarInput,
) (*output.SearchSimilarOutput, error) {
	const op = "comment.UseCase.SearchSimilar"
	logFields := logger.WithFields("operation", op, "destination", i.CommentDestination, "content_length", len(i.Content))

	uc.log.Info("Attempting to search similar comments", logFields()...)

	comments, err := uc.repo.GetByDestination(ctx, i.CommentDestination)
	if err != nil {
		uc.log.Error("Failed to search similar comments", logFields("error", err)...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	validComments := make(model.RawComments, 0, len(comments)/2)
	for _, comment := range comments {
		if strings.Contains(comment.Content, i.Content) {
			if i.Author != "" {
				if comment.Author != i.Author {
					continue
				}
			}
			validComments = append(validComments, comment)
		}
	}

	uc.log.Info("Successfully searched similar comments", logFields("found", len(validComments))...)

	return &output.SearchSimilarOutput{
		Comments: validComments,
		Count:    len(validComments),
	}, nil
}

func (uc *UseCase) GetByIDWithChildren(ctx context.Context, id int) (*model.Comment, error) {
	const op = "comment.UseCase.GetByIDWithChildren"
	logFields := logger.WithFields("operation", op, "comment_id", id)

	uc.log.Info("Attempting to get comment", logFields()...)

	comments, err := uc.repo.GetByIDWithChildren(ctx, id)
	if err != nil {
		uc.log.Error("Failed to get comment", logFields("error", err)...)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tree := service.BuildCommentTree(comments)
	if len(tree) != 1 {
		uc.log.Error("Failed to get comment", logFields(
			"error", errorx.ErrInvalidCommentLength,
			"length", len(tree),
		)...)
		return nil, fmt.Errorf("%s: %w", op, errorx.ErrInvalidCommentLength)
	}

	uc.log.Info("Successfully got comment", logFields()...)

	return tree[0], nil
}
