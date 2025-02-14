package handlers

import (
	"blog/api/helpers"
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/service/article"
	"blog/utils"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/swag/example/celler/httputil"
)

// Article represents the Article handler.
type Article struct {
	ArticleService article.ArticleService
}

// articleInput represents the input for creating or updating an article.
type articleInput struct {
	Title   string `json:"title" binding:"required" example:"My First Article"`                      // Title of the article
	Content string `json:"content" binding:"required" example:"This is the content of the article."` // Content of the article
}

// NewArticleHandler creates a new Article handler.
func NewArticleHandler(articleService article.ArticleService) *Article {
	return &Article{
		ArticleService: articleService,
	}
}

// @Summary Get all articles
// @Description Retrieve all articles
// @Tags articles
// @Produce json
// @Success 200 {object} map[string]interface{} "Successfully retrieved articles"
// @Failure 400 {object} map[string]interface{} "Failed to fetch articles"
// @Router /api/v1/article [get]
func (a *Article) GetAll(ctx *gin.Context) {
	articles, err := a.ArticleService.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Failed to fetch articles",
			map[string]interface{}{
				"error": err.Error(),
				"count": 0,
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Successfully retrieved articles",
		map[string]interface{}{
			"count":     len(articles),
			"timestamp": time.Now(),
			"articles":  articles,
		}))
}

// @Summary Get an article by ID
// @Description Retrieve an article by its ID
// @Tags articles
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} map[string]interface{} "Article retrieved successfully"
// @Failure 404 {object} map[string]interface{} "Article not found"
// @Failure 500 {object} map[string]interface{} "Failed to fetch article"
// @Router /api/v1/article/{id} [get]
func (a *Article) GetByID(ctx *gin.Context) {
	ID := ctx.Param("id")

	article, err := a.ArticleService.GetArticleByID(ctx, uint(helpers.StringToInt(ID)))
	if err != nil {
		if errors.Is(err, postgres_repository.ErrArticleNotFound) {
			ctx.JSON(http.StatusNotFound, helpers.NewHttpResponse(
				http.StatusNotFound,
				"Article not found",
				map[string]interface{}{
					"article_id": ID,
					"error":      err.Error(),
				}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.NewHttpResponse(
			http.StatusInternalServerError,
			"Failed to fetch article",
			map[string]interface{}{
				"article_id": ID,
				"error":      err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Article retrieved successfully",
		map[string]interface{}{
			"article": map[string]interface{}{
				"id":         article.ID,
				"title":      article.Title,
				"content":    article.Content,
				"author":     article.Author,
				"created_at": article.CreatedAt,
				"updated_at": article.UpdatedAt,
			},
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		}))
}


// @Summary Search articles by title
// @Description Search articles by title using a query string
// @Tags articles
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} map[string]interface{} "Articles found successfully"
// @Failure 404 {object} map[string]interface{} "No articles found with given data"
// @Failure 500 {object} map[string]interface{} "Failed to search articles"
// @Router /api/v1/article/search [get]
func (a *Article) Search(ctx *gin.Context) {
	q := ctx.Query("q")

	articles, err := a.ArticleService.GetArticleByTitle(ctx, q)
	if err != nil {
		if errors.Is(err, postgres_repository.ErrArticleNotFound) {
			ctx.JSON(http.StatusNotFound, helpers.NewHttpResponse(
				http.StatusNotFound,
				"No articles found with given data",
				map[string]interface{}{
					"search_query": q,
					"error":        err.Error(),
				}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.NewHttpResponse(
			http.StatusInternalServerError,
			"Failed to search articles",
			map[string]interface{}{
				"search_query": q,
				"error":        err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Articles found successfully",
		map[string]interface{}{
			"articles": articles,
			"count":    len(articles),
			"search_criteria": map[string]interface{}{
				"search_query": q,
			},
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		}))
}


// @Summary Create a new article
// @Description Create a new article with the provided title and content
// @Tags articles
// @Accept json
// @Produce json
// @Param article body articleInput true "Article input"
// @Success 201 {object} map[string]interface{} "Article created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input data"
// @Router /api/v1/article/create [post]
func (a *Article) Create(ctx *gin.Context) {
	var ai articleInput
	err := ctx.ShouldBindJSON(&ai)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
				http.StatusBadRequest,
				"Missing required fields",
				map[string]interface{}{
					"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
					"provided_fields":   ai,
				}))
			return
		}
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid input data",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
				"provided_data":     ai,
			}))
		return
	}

	authorID := uint(ctx.GetInt("id"))
	article, err := a.ArticleService.Create(ctx, ai.Title, ai.Content, authorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Failed to create article",
			map[string]interface{}{
				"error": err.Error(),
				"input": ai,
			}))
		return
	}

	ctx.JSON(http.StatusCreated, helpers.NewHttpResponse(
		http.StatusCreated,
		"Article created successfully",
		map[string]interface{}{
			"article": map[string]interface{}{
				"id":         article.ID,
				"title":      article.Title,
				"content":    article.Content,
				"author_id":  article.AuthorID,
				"created_at": article.CreatedAt,
			},
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
				"status":    "published",
			},
		}))
}

// @Summary Update an article by ID
// @Description Update an article by its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Param article body articleInput true "Article input"
// @Success 200 {object} map[string]interface{} "Article updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid update data"
// @Failure 404 {object} map[string]interface{} "Article not found for update"
// @Router /api/v1/article/{id} [patch]
func (a *Article) UpdateByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var ai = new(articleInput)

	err := ctx.ShouldBind(ai)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
				http.StatusBadRequest,
				"Missing required fields for update",
				map[string]interface{}{
					"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
					"article_id":        id,
					"provided_data":     ai,
				}))
			return
		}

		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid update data",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
				"article_id":        id,
				"provided_data":     ai,
			}))
		return
	}

	article, err := a.ArticleService.Update(ctx,
		uint(helpers.StringToInt(id)),
		ai.Title,
		ai.Content,
	)

	if err != nil {
		if errors.Is(err, postgres_repository.ErrArticleNotFound) {
			ctx.JSON(http.StatusNotFound, helpers.NewHttpResponse(
				http.StatusNotFound,
				"Article not found for update",
				map[string]interface{}{
					"article_id": id,
					"error":      err.Error(),
				}))
			return
		}

		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Failed to update article",
			map[string]interface{}{
				"article_id": id,
				"error":      err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Article updated successfully",
		map[string]interface{}{
			"article": map[string]interface{}{
				"id":         article.ID,
				"title":      article.Title,
				"content":    article.Content,
				"author":     article.Author,
				"updated_at": article.UpdatedAt,
			},
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		}))
}

// @Summary Delete an article by ID
// @Description Delete an article by its ID
// @Tags article
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} map[string]interface{} "Article deleted successfully"
// @Failure 404 {object} map[string]interface{} "Article not found for deletion"
// @Failure 400 {object} map[string]interface{} "Failed to delete article"
// @Router /api/v1/article/{id} [delete]
func (a *Article) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := a.ArticleService.Delete(ctx, uint(helpers.StringToInt(id)))
	if err != nil {
		if errors.Is(err, postgres_repository.ErrArticleNotFound) {
			ctx.JSON(http.StatusNotFound, helpers.NewHttpResponse(
				http.StatusNotFound,
				"Article not found for deletion",
				map[string]interface{}{
					"article_id": id,
					"error":      err.Error(),
				}))
			return
		}
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Failed to delete article",
			map[string]interface{}{
				"article_id": id,
				"error":      err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Article deleted successfully",
		map[string]interface{}{
			"deleted_article_id": id,
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
				"status":    "deleted",
			},
		}))
}
