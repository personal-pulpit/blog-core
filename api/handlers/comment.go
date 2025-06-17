package handlers

import (
	"blog/api/helpers"
	"blog/internal/repository"
	"blog/internal/service/comment"
	"blog/utils"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService comment.CommentService
}

func NewCommentHandler(commentService comment.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// commentInput represents the input for creating a comment.
type createCommentInput struct {
	ArticleID string `json:"article_id" binding:"required" example:"1"`                                // ID of the comment
	Content   string `json:"content" binding:"required" example:"This is the content of the comment."` // Content of the comment
}

type updateComment struct {
	Content string `json:"content" binding:"required"`
}

// @Summary Create a new comment
// @Description Create a new comment with the provided content and comment ID
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body commentInput true "comment input"
// @Success 201 {object} map[string]interface{} "comment created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input data"
// @Router /api/v1/comments [post]
func (a *CommentHandler) Create(ctx *gin.Context) {
	var ci createCommentInput
	err := ctx.ShouldBindJSON(&ci)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			helpers.RespondWithError(ctx, http.StatusBadRequest, "Missing required fields", map[string]interface{}{
				"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
				"provided_fields":   ci,
			})
			return
		}
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid input data", map[string]interface{}{
			"validation_errors": utils.GetValidationError(err),
			"provided_data":     ci,
		})
		return
	}

	userID := uint(ctx.GetInt("id"))
	articleInputID, err := helpers.StringToInt(ci.ArticleID)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid article ID", map[string]interface{}{
			"error":         "Invalid article ID format",
			"provided_data": ci.ArticleID,
		})
		return
	}
	comment, err := a.commentService.AddComment(
		ctx,
		ci.Content,
		userID,
		uint(articleInputID),
	)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Failed to create comment", map[string]interface{}{
			"error": err.Error(),
			"input": ci,
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusCreated, "Comment created successfully", map[string]interface{}{
		"comment": map[string]interface{}{
			"id":         comment.ID,
			"content":    comment.Content,
			"user_id":    comment.UserID,
			"article_id": comment.ArticleID,
			"created_at": comment.CreatedAt,
		},
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
			"status":    "published",
		},
	})
}

// @Summary Update an comment by ID
// @Description Update an comment by its ID
// @Tags comments
// @Accept json
// @Produce json
// @Param id path string true "comment ID"
// @Success 200 {object} map[string]interface{} "comment updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid update data"
// @Failure 404 {object} map[string]interface{} "comment not found for update"
// @Router /api/v1/comments/{id} [patch]
func (a *CommentHandler) UpdateByID(ctx *gin.Context) {
	strID := ctx.Param("id")
	ID, err := helpers.StringToInt(strID)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid comment ID", map[string]interface{}{
			"error":         "Invalid comment ID format",
			"provided_data": strID,
		})
		return
	}

	var input updateComment

	err = ctx.ShouldBindJSON(&input)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid update data", map[string]interface{}{
			"validation_errors": utils.GetValidationError(err),
			"provided_data":     input.Content,
		})
		return
	}

	userID := uint(ctx.GetInt("id"))

	comment, err := a.commentService.UpdateComment(ctx, userID, uint(ID), input.Content)
	if err != nil {
		if errors.Is(err, repository.ErrCommentNotFound) {
			helpers.RespondWithError(ctx, http.StatusNotFound, "Comment not found for update", map[string]interface{}{
				"comment_id": ID,
				"error":      err.Error(),
			})
			return
		}

		helpers.RespondWithError(ctx, http.StatusBadRequest, "Failed to update comment", map[string]interface{}{
			"comment_id": ID,
			"error":      err.Error(),
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Comment updated successfully", map[string]interface{}{
		"comment": map[string]interface{}{
			"id":         comment.ID,
			"content":    comment.Content,
			"user_id":    comment.UserID,
			"article_id": comment.ArticleID,
			"updated_at": comment.UpdatedAt,
		},
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
		},
	})
}

// @Summary Delete an comment by ID
// @Description Delete an comment by its ID
// @Tags comments
// @Produce json
// @Param id path string true "comment ID"
// @Success 200 {object} map[string]interface{} "comment deleted successfully"
// @Failure 404 {object} map[string]interface{} "comment not found for deletion"
// @Failure 400 {object} map[string]interface{} "Failed to delete comment"
// @Router /api/v1/comments/{id} [delete]
func (a *CommentHandler) DeleteByID(ctx *gin.Context) {
	strID := ctx.Param("id")
	ID, err := helpers.StringToInt(strID)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid comment ID", map[string]interface{}{
			"error":         "Invalid comment ID format",
			"provided_data": strID,
		})
		return
	}

	userID := uint(ctx.GetInt("id"))

	err = a.commentService.DeleteComment(ctx, userID, uint(ID))
	if err != nil {
		if errors.Is(err, repository.ErrCommentNotFound) {
			helpers.RespondWithError(ctx, http.StatusNotFound, "Comment not found for deletion", map[string]interface{}{
				"comment_id": ID,
				"error":      err.Error(),
			})
			return
		}
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Failed to delete comment", map[string]interface{}{
			"comment_id": ID,
			"error":      err.Error(),
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Comment deleted successfully", map[string]interface{}{
		"deleted_comment_id": ID,
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
			"status":    "deleted",
		},
	})
}
