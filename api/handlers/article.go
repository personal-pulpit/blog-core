package handlers

import (
	"blog/api/helpers"
	mysql_repository "blog/database/mysql/repo"

	"blog/internal/model"
	"blog/internal/repository"

	"blog/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoIdDetected        = errors.New("no ID detected")
	articleResponseChannel = make(chan helpers.HttpResponse)
)

type (
	Article struct {
		ArticleMysqlRepo repository.ArticleMysqlRepository
		ArticleRedisRepo repository.ArticleRedisRepository
		UserRedisRepo    repository.UserRedisRepository
	}
	ArticleInput struct {
		Title   string `form:"title" binding:"required"`
		Content string `form:"content" binding:"required"`
	}
)

func (a *Article) GetAll(ctx *gin.Context) {
	go func() {
		articles, err := a.ArticleRedisRepo.GetCaches()
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "articles got!", map[string]interface{}{
				"articles": articles,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, articleResponseChannel)
}
func (a *Article) GetByID(ctx *gin.Context) {
	go func() {
		ID := ctx.Param("ID")
		article, err := a.ArticleRedisRepo.GetCacheByID(model.ID(ID))
		if err != nil {
			if errors.Is(err, mysql_repository.ErrArticleNotFound) {
				articleResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), nil)
				return
			}
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		user, err := a.UserRedisRepo.GetCacheByID(article["authorId"])
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "article Got!", map[string]interface{}{
				"title":      article["title"],
				"content":    article["content"],
				"author":     user["username"],
				"created at": article["createdAt"],
				"updated at": article["updatedAt"],
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, articleResponseChannel)
}
func (a *Article) Create(ctx *gin.Context) {
	go func() {
		var ai ArticleInput
		err := ctx.ShouldBind(&ai)
		if err != nil {
			if utils.CheckErrorForWord(err, "required") {
				articleResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest,
					utils.GetValidationError(ErrPleaseCompleteAllFields),
					nil)
				return
			}
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest,
				utils.GetValidationError(err),
				nil)
			return
		}
		//check
		authorId := helpers.GetIdFromToken(ctx)
		article, err := a.ArticleMysqlRepo.Create(
			model.ID(authorId), ai.Title, ai.Content,
		)
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		user, err := a.ArticleRedisRepo.GetCacheByID(model.ID(authorId))
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "article created!", map[string]interface{}{
				"title":      article.Title,
				"content":    article.Content,
				"author":     user["username"],
				"created at": article.CreatedAt,
				"updated at": article.UpdatedAt,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, articleResponseChannel)

}
func (a *Article) UpdateByID(ctx *gin.Context) {
	go func() {
		ID := ctx.Param("ID")
		var ai ArticleInput
		err := ctx.ShouldBind(&ai)
		if err != nil {
			if utils.CheckErrorForWord(err, "required") {
				articleResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest,
					utils.GetValidationError(ErrPleaseCompleteAllFields),
					nil)
				return
			}
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest,
				utils.GetValidationError(err),
				nil)
			return
		}
		article, err := a.ArticleMysqlRepo.UpdateByID(
			ID,
			ai.Title,
			ai.Content,
		)
		if err != nil {
			if errors.Is(err, mysql_repository.ErrArticleNotFound) {
				articleResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), nil)
				return
			}
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		user, err := a.ArticleRedisRepo.GetCacheByID(model.ID(ID))
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "article updated!", map[string]interface{}{
				"title":      article.Title,
				"content":    article.Content,
				"author":     user["username"],
				"created at": article.CreatedAt,
				"updated at": article.UpdatedAt,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, articleResponseChannel)

}
func (a *Article) DeleteByID(ctx *gin.Context) {
	go func() {
		ID := ctx.Param("ID")
		err := a.ArticleMysqlRepo.DeleteByID(ID)
		if err != nil {
			if errors.Is(err, mysql_repository.ErrArticleNotFound) {
				articleResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), nil)
				return
			}
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "article deleted!", map[string]interface{}{},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, articleResponseChannel)

}
