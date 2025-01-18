package auth_middlewares

import (
	"blog/api/helpers"
	"blog/api/helpers/auth_helper"
	"blog/api/helpers/common"
	"blog/config"
	"blog/pkg/auth_manager"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAuthMiddleware struct {
	AuthManager auth_manager.AuthManager
}

func NewUserAuthMiddleware(authManger auth_manager.AuthManager) *UserAuthMiddleware {
	return &UserAuthMiddleware{
		AuthManager: authManger,
	}
}
func (m *UserAuthMiddleware) SetUserStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := auth_helper.GetHeader(ctx, auth_helper.AccessTokenHeader)

		if err != nil {
			ctx.Set("is_logged", false)
		} else {
			if m.AuthManager.IsAccessTokenBlacklisted(ctx, accessToken) {
				ctx.Set("is_logged", false)
				ctx.Next()
			}

			accessTokenClaims, err := m.AuthManager.DecodeAccessToken(ctx, accessToken)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError,
					helpers.NewHttpResponse(
						http.StatusInternalServerError, err.Error(), nil))
				return
			}

			ctx.Set("id", accessTokenClaims.UserID)
			ctx.Set("role", accessTokenClaims.Role)
			ctx.Set("is_logged", true)
		}
		ctx.Next()
	}
}

func (m *UserAuthMiddleware) EnsureLoggedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		is_logged := common.GetUserStatus(ctx)
		if is_logged {
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				helpers.NewHttpResponse(
					http.StatusUnauthorized, ErrYouAreUnAuthorized.Error(), nil))
			return
		}
	}
}

func (m *UserAuthMiddleware) EnsureNotLoggedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		is_logged := common.GetUserStatus(ctx)

		if !is_logged {
			ctx.Next()

		} else {
			port := config.GetConfigInstance().Server.Port
			url := fmt.Sprintf("localhost:%d/api/v1", port)
			ctx.Redirect(http.StatusFound, url)
		}
	}
}

func (m *UserAuthMiddleware) EnsureAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetString("role")
		if role == "admin" {
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				helpers.NewHttpResponse(
					http.StatusUnauthorized, ErrYouAreUnAuthorized.Error(), nil))
			return
		}
	}
}

func (m *UserAuthMiddleware) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		is_logged := common.GetUserStatus(ctx)

		if is_logged {
			accessToken, err := auth_helper.GetHeader(ctx, auth_helper.AccessTokenHeader)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}

			refreshToken, err := auth_helper.GetHeader(ctx, auth_helper.RefreshTokenHeader)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}

			auth_helper.DeleteHeader(ctx, auth_helper.AccessTokenHeader)
			auth_helper.DeleteHeader(ctx, auth_helper.RefreshTokenHeader)

			err = m.AuthManager.DestroyRefreshToken(ctx, refreshToken)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}

			err = m.AuthManager.SetAccessTokenIntoBlacklist(ctx, accessToken, auth_manager.AccessTokenExpr)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest,
				helpers.NewHttpResponse(
					http.StatusBadRequest, ErrSomeTimesWentWrong.Error(), nil))
			return
		}
	}
}
