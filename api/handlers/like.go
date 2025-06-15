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
	ArticleID string `json:"article_id" binding:"required" example:"1"` // ID of the article
}

// @Summary      Like an article
// @Description  Allows the authenticated user to like an article by its ID
// @Tags         Likes
// @Accept       json
// @Produce      json
// @Param        likeInput  body      likeInput  true  "Like input"
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

//@Summary      Get all likes for the authenticated user
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Likes retrieved successfully"
// @Failure      500  {object}  map[string]interface{}  "Failed to retrieve likes"
// @Router       /likes [get]
func (a *LikeHandler) GetUserLikes(ctx *gin.Context) {
	userID := uint(ctx.GetInt("id"))

	likes, err := a.likeService.GetUserLikes(ctx, userID)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve likes", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if len(likes) == 0 {
		helpers.RespondWithSuccess(ctx, http.StatusOK, "No likes found", map[string]interface{}{
			"likes": []interface{}{},
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		})
		return
	}

	var response []map[string]interface{}
	for _, like := range likes {
		response = append(response, map[string]interface{}{
			"id":         like.ID,
			"user_id":    like.UserID,
			"article_id": like.ArticleID,
			"created_at": like.CreatedAt,
			"article": map[string]interface{}{
				"id":    like.Article.ID,
				"title": like.Article.Title,
			},
		})
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Likes retrieved successfully", map[string]interface{}{
		"likes": response,
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
			"count":     len(response),
			"status":    "success",
		},
	})
}

// @Summary      Delete a like
// @Description  Allows the authenticated user to remove their like from an article
// @Tags         Likes
// @Param        id   path      int  true  "Like ID"
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Like deleted successfully"
// @Failure      400  {object}  map[string]interface{}  "Failed to delete like"
// @Failure      404  {object}  map[string]interface{}  "Like not found"
// @Router       /likes/article/{id} [delete]
func (a *LikeHandler) DeleteByID(ctx *gin.Context) {
	articleID := ctx.Param("id")
	userID := uint(ctx.GetInt("id"))

	err := a.likeService.DeleteLike(ctx, uint(helpers.StringToInt(articleID)), userID)
	if err != nil {
		if errors.Is(err, repository.ErrLikeNotFound) {
			helpers.RespondWithError(ctx, http.StatusNotFound, "Like not found for deletion", map[string]interface{}{
				"article_id": articleID,
				"error":      err.Error(),
			})
			return
		}
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Failed to delete like", map[string]interface{}{
			"article_id": articleID,
			"error":      err.Error(),
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Like deleted successfully", map[string]interface{}{
		"deleted_article_like_id": articleID,
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
			"status":    "deleted",
		},
	})
}
