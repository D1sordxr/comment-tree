package params

type Create struct {
	CommentDestination string `json:"comment_destination"`
	Content            string `json:"content" binding:"required"`
	Author             string `json:"author"`
	ParentID           *int   `json:"parent_id"`
}

type GetWithPagination struct {
	Destination string
	CursorID    int
	Limit       int
}
