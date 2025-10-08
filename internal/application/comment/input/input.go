package input

type CreateInput struct {
	CommentDestination string
	Content            string
	Author             string
	ParentID           *int
}

type GetTreeWithPagination struct {
	CommentDestination string
	CursorID           int
}

type GetCommentsWithPagination struct {
	CommentDestination string
	CursorID           int
}

type SearchSimilarInput struct {
	CommentDestination string
	Content            string
	Author             string
}
