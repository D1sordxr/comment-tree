package output

import "github.com/D1sordxr/comment-tree/internal/domain/core/comment/model"

type CreateOutput struct {
	Comment *model.Comment `json:"comment"`
}

type GetTreeOutput struct {
	Tree model.CommentTree `json:"tree"`
}

type GetTreeWithPaginationOutput struct {
	Tree       model.CommentTree `json:"tree"`
	NextCursor int               `json:"next_cursor"`
}

type GetCommentsWithPaginationOutput struct {
	Comments   []model.Comment `json:"comments"`
	NextCursor int             `json:"next_cursor"`
}

type DeleteOutput struct {
	Success bool `json:"success"`
}

type SearchSimilarOutput struct {
	Comments []model.Comment `json:"comments"`
	Count    int             `json:"count"`
}
