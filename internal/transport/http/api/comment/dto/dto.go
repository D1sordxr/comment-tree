package dto

type DefaultResponse map[string]any

type CreateCommentRequest struct {
	ParentID           *int   `json:"parent_id"`
	Content            string `json:"content" binding:"required" validate:"required,max=255"`
	Author             string `json:"author" binding:"required" validate:"required,max=255"`
	CommentDestination string `json:"comment_destination" binding:"required" validate:"required,max=255"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
