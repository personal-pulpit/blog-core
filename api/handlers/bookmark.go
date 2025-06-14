package handlers

import (
	"blog/api/helpers"
	"blog/internal/repository"
	"blog/internal/service/bookmark"
	"blog/utils"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct {
	bookmarkService bookmark.BookmarkService
}

func NewBookmarkHandler(bookmarkService bookmark.BookmarkService) *BookmarkHandler {
	return &BookmarkHandler{
		bookmarkService: bookmarkService,
	}
}

type bookmarkInput struct {
	ArticleID string `json:"article_id" binding:"required" example:"1"`
}

// @Summary      Create a bookmark
// @Description  Allows the authenticated user to create a bookmark for an article by its ID
// @Tags         Bookmarks
// @Accept       json
// @Produce      json
// @Param        bookmarkInput  body      bookmarkInput  true  "Bookmark input"
// @Security     BearerAuth
// @Success      201  {object}  map[string]interface{}  "Bookmark created successfully"
// @Failure      400  {object}  map[string]interface{}  "Invalid input or article ID"
// @Router       /bookmarks [post]
func (h *BookmarkHandler) Create(ctx *gin.Context) {
	var bi bookmarkInput
	err := ctx.ShouldBindJSON(&bi)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			helpers.RespondWithError(ctx, http.StatusBadRequest, "Missing required fields", map[string]interface{}{
				"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
				"provided_fields":   bi,
			})
			return
		}
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid input data", map[string]interface{}{
			"validation_errors": utils.GetValidationError(err),
			"provided_data":     bi,
		})
		return
	}

	userID := uint(ctx.GetInt("id"))
	bookmark, err := h.bookmarkService.CreateBookmark(
		ctx,
		userID,
		uint(helpers.StringToInt(bi.ArticleID)),
	)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Failed to create bookmark", map[string]interface{}{
			"error": err.Error(),
			"input": bi,
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusCreated, "Bookmark created successfully", map[string]interface{}{
		"bookmark": map[string]interface{}{
			"id":         bookmark.ID,
			"user_id":    bookmark.UserID,
			"article_id": bookmark.ArticleID,
			"created_at": bookmark.CreatedAt,
		},

		"metadata": map[string]interface{}{
			"timestamp": time.Now()},
	})
}


// @Summary      Get all bookmarks for the authenticated user
// @Description  Retrieve all bookmarks for the authenticated user
// @Tags         Bookmarks
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "Bookmarks retrieved successfully"
// @Failure      404  {object}  map[string]interface{}  "No bookmarks found"
// @Failure      500  {object}  map[string]interface{}  "Failed to retrieve bookmarks"
// @Router       /bookmarks [get]
func (a *BookmarkHandler) GetUsersBookmarks(ctx *gin.Context) {
	userID := uint(ctx.GetInt("id"))

	bookmarks, err := a.bookmarkService.GetUsersBookmarks(ctx, userID)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve bookmarks", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if len(bookmarks) == 0 {
		helpers.RespondWithSuccess(ctx, http.StatusOK, "No bookmarks found", map[string]interface{}{
			"bookmarks": []interface{}{},
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		})
		return
	}

	var response []map[string]interface{}
	for _, bookmark := range bookmarks {
		response = append(response, map[string]interface{}{
			"id":         bookmark.ID,
			"user_id":    bookmark.UserID,
			"article_id": bookmark.ArticleID,
			"created_at": bookmark.CreatedAt,
			"article": map[string]interface{}{
				"id":    bookmark.Article.ID,
				"title": bookmark.Article.Title,
			},
		})
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Bookmarks retrieved successfully", map[string]interface{}{
		"bookmarks": response,
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
			"count":     len(response),
			"status":    "success",
		},
	})
}

// @Summary      Get a bookmark by ID
// @Description  Retrieve a bookmark by its ID for the authenticated user
// @Tags         Bookmarks
// @Produce      json
// @Param        id  path      string  true  "Bookmark ID" example:"1"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "Bookmark retrieved successfully"
// @Failure      404  {object}  map[string]interface{}  "Bookmark not found"
// @Failure      400  {object}  map[string]interface{}  "Invalid bookmark ID"
func (a *BookmarkHandler) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := uint(ctx.GetInt("id"))

	err := a.bookmarkService.DeleteBookmark(ctx, userID, uint(helpers.StringToInt(id)))
	if err != nil {
		if errors.Is(err, repository.ErrLikeNotFound) {
			helpers.RespondWithError(ctx, http.StatusNotFound, "Like not found for deletion", map[string]interface{}{
				"like_id": id,
				"error":   err.Error(),
			})
			return
		}
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Failed to delete bookmark", map[string]interface{}{
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
