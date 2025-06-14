package handlers

import (
	"blog/api/helpers"
	"blog/internal/repository"
	"blog/internal/service/like"
	"blog/utils"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeService like.LikeService
}

func NewLikeHandler(likeService like.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
	}
}

// likeInput represents the input for creating a like.
type likeInput struct {
	ArticleID string `json:"article_id" binding:"required" example:"1"`                             // ID of the article
}

// @Summary      Like an article
// @Description  Allows the authenticated user to like an article by its ID
// @Tags         Likes
// @Accept       json
// @Produce      json
// @Param        likeInput  body      likeInput  true  "Like input"
// @Security     BearerAuth
// @Success      201  {object}  map[string]interface{}  "Like created successfully"
// @Failure      400  {object}  map[string]interface{}  "Invalid input or article ID"
// @Router       /likes [post]
func (a *LikeHandler) Create(ctx *gin.Context) {
	var li likeInput
	err := ctx.ShouldBindJSON(&li)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			helpers.RespondWithError(ctx, http.StatusBadRequest, "Missing required fields", map[string]interface{}{
				"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
				"provided_fields":   li,
			})
			return
		}
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid input data", map[string]interface{}{
			"validation_errors": utils.GetValidationError(err),
			"provided_data":     li,
		})
		return
	}

	userID := uint(ctx.GetInt("id"))
	like, err := a.likeService.CreateLike(
		ctx,
		userID,
		uint(helpers.StringToInt(li.ArticleID)),
	)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Failed to create like", map[string]interface{}{
			"error": err.Error(),
			"input": li,
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusCreated, "Like created successfully", map[string]interface{}{
		"like": map[string]interface{}{
			"id":         like.ID,
			"user_id":    like.UserID,
			"article_id": like.ArticleID,
			"created_at": like.CreatedAt,
		},

		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
			"status":    "published",
		},
	})
}


// @Summary      Delete a like
// @Description  Allows the authenticated user to remove their like from an article
// @Tags         Likes
// @Param        id   path      int  true  "Like ID"
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "Like deleted successfully"
// @Failure      400  {object}  map[string]interface{}  "Failed to delete like"
// @Failure      404  {object}  map[string]interface{}  "Like not found"
// @Router       /likes/{id} [delete]
func (a *LikeHandler) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := uint(ctx.GetInt("id"))

	err := a.likeService.DeleteLike(ctx, userID, uint(helpers.StringToInt(id)))
	if err != nil {
		if errors.Is(err, repository.ErrLikeNotFound) {
			helpers.RespondWithError(ctx, http.StatusNotFound, "Like not found for deletion", map[string]interface{}{
				"like_id": id,
				"error":   err.Error(),
			})
			return
		}
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Failed to delete like", map[string]interface{}{
			"like_id": id,
			"error":   err.Error(),
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Like deleted successfully", map[string]interface{}{
		"deleted_like_id": id,
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
			"status":    "deleted",
		},
	})
}
