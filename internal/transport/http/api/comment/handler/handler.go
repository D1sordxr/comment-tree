package handler

import (
	"github.com/D1sordxr/comment-tree/internal/application/comment/input"
	"github.com/D1sordxr/comment-tree/internal/application/comment/port"
	"github.com/D1sordxr/comment-tree/internal/domain/core/shared/validator"
	"github.com/D1sordxr/comment-tree/internal/transport/http/api/comment/dto"
	"github.com/D1sordxr/comment-tree/pkg/httputil"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"strconv"
)

type Handler struct {
	uc port.UseCase
}

func New(uc port.UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) CreateComment(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := httputil.BindWithCustomValidation(
		c,
		&req,
		validator.ValidateStruct, // just an experiment
	); err != nil {
		c.JSON(400, dto.DefaultResponse{"error": "Invalid request"})
		return
	}

	comment, err := h.uc.Create(c.Request.Context(), input.CreateInput{
		ParentID:           req.ParentID,
		Content:            req.Content,
		Author:             req.Author,
		CommentDestination: req.CommentDestination,
	})
	if err != nil {
		c.JSON(500, dto.DefaultResponse{"error": err.Error()})
		return
	}

	c.JSON(201, dto.DefaultResponse{"data": comment})
}

func (h *Handler) GetComment(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, dto.DefaultResponse{"error": "Invalid comment ID"})
		return
	}

	comment, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, dto.DefaultResponse{"error": err.Error()})
		return
	}

	if comment == nil {
		c.JSON(404, dto.DefaultResponse{"error": "Comment not found"})
		return
	}

	c.JSON(200, dto.DefaultResponse{"data": comment})
}

func (h *Handler) GetComments(c *ginext.Context) {
	parentID, _ := strconv.Atoi(c.Query("parent"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sort := c.DefaultQuery("sort", "created_at_desc")
	search := c.Query("search")

	var parentIDPtr *int
	if c.Query("parent") != "" {
		parentIDPtr = &parentID
	}

	comments, err := h.uc.GetMany(c.Request.Context(), input.GetManyInput{
		ParentID: parentIDPtr,
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Search:   search,
	})
	if err != nil {
		c.JSON(500, dto.DefaultResponse{"error": err.Error()})
		return
	}

	response := dto.CommentsResponse{
		Comments: comments,
		Total:    len(comments),
		Page:     page,
		PageSize: pageSize,
	}

	c.JSON(200, dto.DefaultResponse{"data": response})
}

func (h *Handler) DeleteComment(c *ginext.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, dto.DefaultResponse{"error": "Invalid comment ID"})
		return
	}

	comment, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, dto.DefaultResponse{"error": err.Error()})
		return
	}

	if comment == nil {
		c.JSON(404, dto.DefaultResponse{"error": "Comment not found"})
		return
	}

	err = h.uc.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, dto.DefaultResponse{"error": err.Error()})
		return
	}

	c.JSON(200, dto.DefaultResponse{"message": "Comment deleted successfully"})
}

func (h *Handler) SearchComments(c *ginext.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if query == "" {
		c.JSON(400, dto.DefaultResponse{"error": "Search query is required"})
		return
	}

	comments, err := h.uc.GetMany(c.Request.Context(), input.GetManyInput{
		Search:   query,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.JSON(500, dto.DefaultResponse{"error": err.Error()})
		return
	}

	response := dto.CommentsResponse{
		Comments: comments,
		Total:    len(comments),
		Page:     page,
		PageSize: pageSize,
	}

	c.JSON(200, dto.DefaultResponse{
		"data": response,
		"meta": map[string]interface{}{
			"search_query": query,
		},
	})
}

func (h *Handler) RegisterRoutes(router *ginext.RouterGroup) {
	router.POST("/comments", h.CreateComment)
	router.GET("/comments/:id", h.GetComment)
	router.GET("/comments", h.GetComments)
	router.DELETE("/comments/:id", h.DeleteComment)
	router.GET("/comments/search", h.SearchComments)
}
