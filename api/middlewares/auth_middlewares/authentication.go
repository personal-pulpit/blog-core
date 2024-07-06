package middlewares

import (
	"blog/api/helpers"
	"blog/api/helpers/auth_helper"
	"blog/api/helpers/common"
	"blog/pkg/auth_manager"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAuthMiddleware struct {
	AuthManager auth_manager.AuthManager
	AuthHelper  auth_helper.AuthHeaderHelper
}

func NewUserAuthMiddelware(authManger auth_manager.AuthManager, authHelper auth_helper.AuthHeaderHelper) *UserAuthMiddleware {
	return &UserAuthMiddleware{
		AuthManager: authManger,
		AuthHelper:  authHelper,
	}
}
func (m *UserAuthMiddleware) SetUserStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := m.AuthHelper.GetHeader(ctx, auth_helper.AccessTokenHeader)

		if err != nil {
			ctx.Set("is_logged", false)
		} else {
			cliams, err := m.AuthManager.DecodeToken(accessToken, auth_manager.AccessToken)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError,
					helpers.NewHttpResponse(
						http.StatusInternalServerError, err.Error(), nil))
				return
			}

			ctx.Set("id", cliams.ID)
			ctx.Set("role", cliams.Role)
			ctx.Set("is_logged", true)
			ctx.Set("is_admin", common.IsAdmin(cliams.Role))

			ctx.Next()
		}

		verifyEmailToken, err := m.AuthHelper.GetHeader(ctx, auth_helper.VerifyEmailTokenHeader)

		if err != nil {
			ctx.Next()
		} else {
			cliams, err := m.AuthManager.DecodeToken(verifyEmailToken, auth_manager.VerifyEmail)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError,
					helpers.NewHttpResponse(
						http.StatusInternalServerError, err.Error(), nil))
				return
			}

			ctx.Set("id", cliams.ID)

			ctx.Next()
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
			accsessToken, err := m.AuthHelper.GetHeader(ctx, auth_helper.AccessTokenHeader)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}

			RefreshToken, err := m.AuthHelper.GetHeader(ctx, auth_helper.RefreshTokenHeader)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}

			m.AuthHelper.DeleteHeader(ctx, auth_helper.AccessTokenHeader)

			m.AuthHelper.DeleteHeader(ctx, auth_helper.RefreshTokenHeader)

			err = m.AuthManager.Destroy(accsessToken)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}

			err = m.AuthManager.Destroy(RefreshToken)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}

		}
	}
}
func (m *UserAuthMiddleware) EnsureAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetBool("is_admin") {
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
			accsessToken, err := m.AuthHelper.GetHeader(ctx, auth_helper.AccessTokenHeader)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}

			RefreshToken, err := m.AuthHelper.GetHeader(ctx, auth_helper.RefreshTokenHeader)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}

			m.AuthHelper.DeleteHeader(ctx, auth_helper.AccessTokenHeader)

			m.AuthHelper.DeleteHeader(ctx, auth_helper.RefreshTokenHeader)

			err = m.AuthManager.Destroy(accsessToken)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}

			err = m.AuthManager.Destroy(RefreshToken)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest,
					helpers.NewHttpResponse(
						http.StatusBadRequest, err.Error(), nil))
				return
			}

		} else {
			ctx.AbortWithStatusJSON(http.StatusBadRequest,
				helpers.NewHttpResponse(
					http.StatusBadRequest, ErrSomeTimesWentWrong.Error(), nil))
			return
		}
	}
}
