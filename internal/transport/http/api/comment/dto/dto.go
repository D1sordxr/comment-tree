package dto

import (
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/model"
)

type DefaultResponse map[string]any

type CreateCommentRequest struct {
	ParentID           *int   `json:"parent_id"`
	Content            string `json:"content" binding:"required" validate:"required,max=255"`
	Author             string `json:"author" binding:"required" validate:"required,max=255"`
	CommentDestination string `json:"comment_destination" binding:"required" validate:"required,max=255"`
}

type CommentsResponse struct {
	Comments []model.Comment `json:"comments"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}
