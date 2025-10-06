package converters

import (
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/params"
	"github.com/D1sordxr/comment-tree/internal/infrastructure/storage/postgres/repositories/comment/gen"
	"github.com/D1sordxr/comment-tree/pkg/sqlutil"
)

func ConvertCreateParams(p params.Create) gen.CreateCommentParams {
	createParams := gen.CreateCommentParams{
		CommentDestination: p.CommentDestination,
		Content:            p.Content,
		Author:             sqlutil.ToNullString(p.Author),
	}

	if p.ParentID != nil {
		parentID := int32(*p.ParentID)
		createParams.ParentID = sqlutil.ToNullInt32(parentID)
	} else {
		createParams.ParentID = sqlutil.ToNullInt32(0)
	}

	return createParams
}

func ConvertGetRootIDsWithPaginationParams(p params.GetWithPagination) gen.GetRootCommentIDsWithPaginationParams {
	return gen.GetRootCommentIDsWithPaginationParams{
		CommentDestination: p.Destination,
		ID:                 int32(p.CursorID),
		Limit:              int32(p.Limit),
	}
}
