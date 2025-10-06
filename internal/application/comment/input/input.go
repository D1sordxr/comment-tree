package input

type CreateInput struct {
	CommentDestination string `json:"comment_destination" binding:"required"`
	Content            string `json:"content" binding:"required"`
	Author             string `json:"author"`
	ParentID           *int   `json:"parent_id"`
}

type SearchSimilarInput struct {
	CommentDestination string `json:"comment_destination" binding:"required"`
	Content            string `json:"content" binding:"required"`
	Author             string `json:"author"`
}

type GetTreeWithPagination struct {
	CursorID           int    `json:"cursor"`
	CommentDestination string `json:"comment_destination"`
}
