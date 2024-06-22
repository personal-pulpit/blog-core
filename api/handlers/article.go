package handlers

import (
	"blog/api/helpers"
	mysql_repository "blog/database/mysql_repo"
	"blog/internal/repository"

	"blog/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoIdDetected        = errors.New("no ID detected")
	articleResponseChannel = make(chan helpers.HttpResponse)
)

type (
	Article struct {
		ArticleRepo repository.ArticleRepository
		UserRepo    repository.UserRepository
	}
	ArticleInput struct {
		Title   string `form:"title" binding:"required"`
		Content string `form:"content" binding:"required"`
	}
)

func (a *Article) GetAll(ctx *gin.Context) {
	go func() {
		articles, err := a.ArticleRepo.GetAll()
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
		article, err := a.ArticleRepo.GetByID(ID)
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
		username, err := a.UserRepo.GetUsernameById(article["authorId"])
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "article Got!", map[string]interface{}{
				"title":      article["title"],
				"content":    article["content"],
				"author":     username,
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
		authorid := helpers.GetIdFromToken(ctx)
		article, err := a.ArticleRepo.Create(
			authorid, ai.Title, ai.Content,
		)
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		username, err := a.UserRepo.GetUsernameById(authorid)
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "article created!", map[string]interface{}{
				"title":      article.Title,
				"content":    article.Content,
				"author":     username,
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
		article, err := a.ArticleRepo.UpdateByID(
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
		username, err := a.UserRepo.GetUsernameById(strconv.Itoa(int(article.AuthorId)))
		if err != nil {
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		articleResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "article updated!", map[string]interface{}{
				"title":      article.Title,
				"content":    article.Content,
				"author":     username,
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
		err := a.ArticleRepo.DeleteByID(ID)
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
