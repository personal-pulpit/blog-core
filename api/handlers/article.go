package handlers

import (
	"blog/api/helpers"
	postgres_repository "blog/database/postgres/repo"

	"blog/internal/service/article"
	"blog/internal/service/user"

	"blog/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoIdDetected        = errors.New("no id detected")
	articleResponseChannel = make(chan helpers.HttpResponse)
)

type (
	Article struct {
		UserService    user.UserService
		ArticleService article.ArticleService
	}
	ArticleInput struct {
		Title   string `form:"title" binding:"required"`
		Content string `form:"content" binding:"required"`
	}
)

func (a *Article) GetAll(ctx *gin.Context) {
	go func() {
		articles, err := a.ArticleService.GetAll()
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
func (a *Article) GetById(ctx *gin.Context) {
	go func() {
		id := ctx.Param("id")
		article, err := a.ArticleService.GetArticleById(id)
		if err != nil {
			if errors.Is(err, postgres_repository.ErrArticleNotFound) {
				articleResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), nil)
				return
			}
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		user, err := a.UserService.GetUserProfile(article.AuthorId)
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "article Got!", map[string]interface{}{
				"title":      article.Title,
				"content":    article.Content,
				"author":     user.FirstName + " " + user.LastName,
				"created at": article.CreatedAt,
				"updated at": article.UpdatedAt,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, articleResponseChannel)
}

func (a *Article) GetByTitle(ctx *gin.Context) {
	go func() {
		title := ctx.Query("title")

		article, err := a.ArticleService.GetArticleByTitle(title)

		if err != nil {
			if errors.Is(err, postgres_repository.ErrArticleNotFound) {
				articleResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), nil)
				return
			}
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}

		user, err := a.UserService.GetUserProfile(article.AuthorId)
		
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "article Got!", map[string]interface{}{
				"title":      article.Title,
				"content":    article.Content,
				"author":     user.FirstName + " " + user.LastName,
				"created at": article.CreatedAt,
				"updated at": article.UpdatedAt,
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
		authorId := ctx.GetString("id")
		article, err := a.ArticleService.Create(
			authorId, ai.Title, ai.Content,
		)
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		user, err := a.UserService.GetUserProfile(authorId)
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "article created!", map[string]interface{}{
				"title":      article.Title,
				"content":    article.Content,
				"author":     user.FirstName + " " + user.LastName,
				"created at": article.CreatedAt,
				"updated at": article.UpdatedAt,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, articleResponseChannel)

}
func (a *Article) UpdateById(ctx *gin.Context) {
	go func() {
		id := ctx.Param("id")
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
		article, err := a.ArticleService.Update(
			id,
			ai.Title,
			ai.Content,
		)
		if err != nil {
			if errors.Is(err, postgres_repository.ErrArticleNotFound) {
				articleResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), nil)
				return
			}
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		user, err := a.UserService.GetUserProfile(article.AuthorId)
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "article updated!", map[string]interface{}{
				"title":      article.Title,
				"content":    article.Content,
				"author":     user.ID,
				"created at": article.CreatedAt,
				"updated at": article.UpdatedAt,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, articleResponseChannel)

}
func (a *Article) DeleteById(ctx *gin.Context) {
	go func() {
		id := ctx.Param("id")
		err := a.ArticleService.Delete(id)
		if err != nil {
			if errors.Is(err, postgres_repository.ErrArticleNotFound) {
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
