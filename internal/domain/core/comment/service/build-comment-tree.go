package service

import (
	"github.com/D1sordxr/comment-tree/internal/domain/core/comment/model"
	"sort"
)

func BuildCommentTree(comments model.RawComments) model.CommentTree {
	commentMap := make(map[int]*model.Comment, len(comments))

	for i := range comments {
		comment := &comments[i]
		comment.Children = make([]*model.Comment, 0)
		commentMap[comment.ID] = comment
	}

	var roots []*model.Comment
	for _, comment := range commentMap {
		if comment.ParentID == nil {
			roots = append(roots, comment)
		} else {
			if parent, exists := commentMap[*comment.ParentID]; exists {
				parent.Children = append(parent.Children, comment)
			}
		}
	}

	sortComments(roots)
	return roots
}

func sortComments(comments model.CommentTree) {
	sort.Slice(comments, func(i, j int) bool {
		return comments[i].CreatedAt.Before(comments[j].CreatedAt)
	})

	for _, comment := range comments {
		if len(comment.Children) > 0 {
			sortComments(comment.Children)
		}
	}
}
