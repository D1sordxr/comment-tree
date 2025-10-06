package model

import "time"

type Comment struct {
	ID                 int         `json:"id"`
	ParentID           *int        `json:"parent_id,omitempty"`
	CommentDestination string      `json:"comment_destination"`
	Author             string      `json:"author"`
	Content            string      `json:"content"`
	Children           CommentTree `json:"children,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at,omitempty"`
}

type (
	RawComments []Comment
	CommentTree []*Comment
)
