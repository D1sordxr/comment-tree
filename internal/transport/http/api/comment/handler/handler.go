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
		validator.ValidateStruct,
	); err != nil {
		c.JSON(400, dto.ErrorResponse{
			Error:   "Invalid request",
			Details: err.Error(),
		})
		return
	}

	result, err := h.uc.Create(c.Request.Context(), input.CreateInput{
		ParentID:           req.ParentID,
		Content:            req.Content,
		Author:             req.Author,
		CommentDestination: req.CommentDestination,
	})
	if err != nil {
		c.JSON(500, dto.ErrorResponse{
			Error:   "Failed to create comment",
			Details: err.Error(),
		})
		return
	}

	c.JSON(201, dto.SuccessResponse{
		Message: "Comment created successfully",
		Data:    result,
	})
}

func (h *Handler) GetCommentTree(c *gin.Context) {
	destination := c.Query("destination")
	cursor, _ := strconv.Atoi(c.DefaultQuery("cursor", "0"))

	if destination == "" {
		c.JSON(400, dto.ErrorResponse{
			Error: "Destination parameter is required",
		})
		return
	}

	var (
		result interface{}
		err    error
	)

	if cursor > 0 {
		result, err = h.uc.GetTreeWithPagination(c.Request.Context(), input.GetTreeWithPagination{
			CommentDestination: destination,
			CursorID:           cursor,
		})
	} else {
		result, err = h.uc.GetTreeByDestination(c.Request.Context(), destination)
	}

	if err != nil {
		c.JSON(500, dto.ErrorResponse{
			Error:   "Failed to get comments",
			Details: err.Error(),
		})
		return
	}

	c.JSON(200, dto.SuccessResponse{
		Data: result,
	})
}

func (h *Handler) GetComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, dto.ErrorResponse{
			Error: "Invalid comment ID",
		})
		return
	}

	treeResult, err := h.uc.GetByIDWithChildren(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, dto.ErrorResponse{
			Error:   "Failed to get comment",
			Details: err.Error(),
		})
		return
	}

	c.JSON(200, dto.SuccessResponse{
		Data: treeResult,
	})
}

func (h *Handler) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, dto.ErrorResponse{
			Error: "Invalid comment ID",
		})
		return
	}

	result, err := h.uc.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, dto.ErrorResponse{
			Error:   "Failed to delete comment",
			Details: err.Error(),
		})
		return
	}

	if !result.Success {
		c.JSON(500, dto.ErrorResponse{
			Error: "Failed to delete comment",
		})
		return
	}

	c.JSON(200, dto.SuccessResponse{
		Message: "Comment deleted successfully",
		Data:    result,
	})
}

func (h *Handler) SearchComments(c *gin.Context) {
	query := c.Query("q")
	destination := c.Query("destination")
	author := c.Query("author")

	if query == "" {
		c.JSON(400, dto.ErrorResponse{
			Error: "Search query is required",
		})
		return
	}

	if destination == "" {
		c.JSON(400, dto.ErrorResponse{
			Error: "Destination parameter is required for search",
		})
		return
	}

	result, err := h.uc.SearchSimilar(c.Request.Context(), input.SearchSimilarInput{
		CommentDestination: destination,
		Content:            query,
		Author:             author,
	})
	if err != nil {
		c.JSON(500, dto.ErrorResponse{
			Error:   "Search failed",
			Details: err.Error(),
		})
		return
	}

	c.JSON(200, dto.SuccessResponse{
		Data: result,
	})
}

func (h *Handler) RegisterRoutes(router *ginext.RouterGroup) {
	// – POST /comments — создание комментария (с указанием родительского)
	router.POST("/comments", h.CreateComment)

	// – GET /comments?parent={id} — получение комментария и всех вложенных
	router.GET("/comments", h.GetComment)

	// – DELETE /comments/{id} — удаление комментария и всех вложенных под ним
	router.DELETE("/comments/:id", h.DeleteComment)

	// Additional:
	// – GET /comments/tree?destination=post_123&cursor=0 – получение комментариев с пагинацией
	router.GET("/comments/tree", h.GetCommentTree)

	// GET /comments/search?destination=post_123&q=search_term – полнотекстовый поиск по комментариям
	router.GET("/comments/search", h.SearchComments)
}
