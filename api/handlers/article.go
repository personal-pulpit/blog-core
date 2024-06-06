package handlers

import (
	"blog/api/helpers"
	"blog/pkg/data/repo"
	db "blog/pkg/data/repo/DB"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoIdDetected        = errors.New("no id detected")
	articleResponseChannel = make(chan helpers.HttpResponse)
)

type (
	Article struct {
		ArticleRepo repo.ArticleDB
		UserRepo    repo.UserDB
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
func (a *Article) GetById(ctx *gin.Context) {
	go func() {
		id := ctx.Param("id")
		article, err := a.ArticleRepo.GetById(id)
		if err != nil {
			if errors.Is(err, db.ErrArticleNotFound) {
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
			articleResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
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
func (a *Article) UpdateById(ctx *gin.Context) {
	go func() {
		id := ctx.Param("id")
		var ai ArticleInput
		err := ctx.ShouldBind(&ai)
		if err != nil {
				articleResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), nil)
				return
		}
		article, err := a.ArticleRepo.UpdateById(
			id,
			ai.Title,
			ai.Content,
		)
		if err != nil {
			if errors.Is(err, db.ErrArticleNotFound) {
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
func (a *Article) DeleteById(ctx *gin.Context) {
	go func() {
		id := ctx.Param("id")
		err := a.ArticleRepo.DeleteById(id)
		if err != nil {
			if errors.Is(err, db.ErrArticleNotFound) {
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
