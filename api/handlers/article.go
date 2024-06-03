package handlers

import (
	"blog/api/helpers"
	"blog/pkg/data/repo"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoIdDetected = errors.New("no id detected")
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
	articles, err := a.ArticleRepo.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in getting articles!", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
		http.StatusOK, "articles got!", map[string]interface{}{
			"articles": articles,
		},
	))
}
func (a *Article) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	article, err := a.ArticleRepo.GetById(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in getting article!", err),
		)
		return
	}
	username, err := a.UserRepo.GetUsernameById(article["authorId"])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in creatig article", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
		http.StatusOK, "article Got!", map[string]interface{}{
			"title":      article["title"],
			"content":    article["content"],
			"author":     username,
			"created at": article["createdAt"],
			"updated at": article["updatedAt"],
		},
	))
}
func (a *Article) Create(ctx *gin.Context) {
	var ai ArticleInput
	err := ctx.ShouldBind(&ai)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "sometimes went wrong", err),
		)
		return
	}
	//check
	authorid := helpers.GetIdFromToken(ctx)
	article, err := a.ArticleRepo.Create(
		authorid, ai.Title, ai.Content,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in creatig article", err),
		)
		return
	}
	username, err := a.UserRepo.GetUsernameById(authorid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in creatig article", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
		http.StatusOK, "article created!", map[string]interface{}{
			"title":      article.Title,
			"content":    article.Content,
			"author":     username,
			"created at": article.CreatedAt,
			"updated at": article.UpdatedAt,
		},
	))
}
func (a *Article) UpdateById(ctx *gin.Context) {
	id, exits := ctx.GetQuery("id")
	if !exits {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in modifying article", ErrNoIdDetected),
		)
		return
	}
	var ai ArticleInput
	err := ctx.ShouldBind(&ai)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "sometimes went wrong", err),
		)
		return
	}
	article, err := a.ArticleRepo.UpdateById(
		id,
		ai.Title,
		ai.Content,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in updating article!", err),
		)
		return
	}
	username, err := a.UserRepo.GetUsernameById(strconv.Itoa(int(article.AuthorId)))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in creatig article", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
		http.StatusOK, "article updated!", map[string]interface{}{
			"title":      article.Title,
			"content":    article.Content,
			"author":     username,
			"created at": article.CreatedAt,
			"updated at": article.UpdatedAt,
		},
	))
}
func (a *Article) DeleteById(ctx *gin.Context) {
	id := ctx.Param("id")
	err := a.ArticleRepo.DeleteById(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in deleting article!", err))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
		http.StatusOK, "article deleted!", map[string]interface{}{},
	))
}
