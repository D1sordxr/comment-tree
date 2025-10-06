package converters

import (
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/model"
	"github.com/D1sordxr/comment-tree/internal/infrastructure/storage/postgres/repositories/comment/gen"
	"time"
)

func ConvertGenToDomain(raw gen.Comment) model.Comment {
	comment := model.Comment{
		ID:                 int(raw.ID),
		CommentDestination: raw.CommentDestination,
		Content:            raw.Content,
		Children:           make(model.CommentTree, 0),
	}

	if raw.ParentID.Valid {
		parentID := int(raw.ParentID.Int32)
		comment.ParentID = &parentID
	}

	if raw.Author.Valid {
		comment.Author = raw.Author.String
	} else {
		comment.Author = "Anonymous"
	}

	if raw.CreatedAt.Valid {
		comment.CreatedAt = raw.CreatedAt.Time
	} else {
		comment.CreatedAt = time.Now()
	}

	if raw.UpdatedAt.Valid {
		comment.UpdatedAt = raw.UpdatedAt.Time
	}

	return comment
}

func ConvertGenSliceToDomain(rawComments []gen.Comment) []model.Comment {
	result := make([]model.Comment, len(rawComments))
	for i, raw := range rawComments {
		result[i] = ConvertGenToDomain(raw)
	}
	return result
}

func ConvertGetCommentsWithChildrenRowToModel(raw gen.GetCommentsWithChildrenRow) model.Comment {
	comment := model.Comment{
		ID:                 int(raw.ID),
		CommentDestination: raw.CommentDestination,
		Content:            raw.Content,
		Children:           make(model.CommentTree, 0),
	}

	if raw.ParentID.Valid {
		parentID := int(raw.ParentID.Int32)
		comment.ParentID = &parentID
	}

	if raw.Author.Valid {
		comment.Author = raw.Author.String
	} else {
		comment.Author = "Anonymous"
	}

	if raw.CreatedAt.Valid {
		comment.CreatedAt = raw.CreatedAt.Time
	} else {
		comment.CreatedAt = time.Now()
	}

	if raw.UpdatedAt.Valid {
		comment.UpdatedAt = raw.UpdatedAt.Time
	}

	return comment
}
