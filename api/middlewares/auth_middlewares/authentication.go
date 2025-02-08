package auth_middlewares

import (
	"blog/api/helpers"
	"blog/api/helpers/auth_helper"
	"blog/api/helpers/common"
	"blog/internal/model"
	"blog/pkg/auth_manager"
	"net/http"
	"strconv"
	"time"

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
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewHttpResponse(
					http.StatusUnauthorized,
					"Token Blacklisted",
					map[string]interface{}{
						"message":    "Your session token has been invalidated",
						"suggestion": "Please login again to obtain a new token",
						"metadata": map[string]interface{}{
							"timestamp":    time.Now(),
							"token_status": "blacklisted",
						},
					}))
				return
			}

			accessTokenClaims, err := m.AuthManager.DecodeAccessToken(ctx, accessToken)
			if err != nil {
				if ctx.Request.URL.Path == "/api/v1/auth/refresh-token" {
					ctx.Set("is_logged", false)
					return
				}

				ctx.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewHttpResponse(
					http.StatusUnauthorized,
					"Invalid Token",
					map[string]interface{}{
						"message":    "Unable to verify your authentication token",
						"suggestion": "Please ensure your token is valid or login again",
						"error":      err.Error(),
						"metadata": map[string]interface{}{
							"timestamp":    time.Now(),
							"token_status": "invalid",
						},
					}))
				return
			}

			ctx.Set("id", helpers.StringToInt(accessTokenClaims.UserID))
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
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, helpers.NewHttpResponse(
				http.StatusUnauthorized,
				"Authentication Required",
				map[string]interface{}{
					"message":    "You must be logged in to access this resource",
					"suggestion": "Please login to continue",
					"metadata": map[string]interface{}{
						"timestamp":      time.Now(),
						"request_path":   ctx.Request.URL.Path,
						"request_method": ctx.Request.Method,
					},
				}))
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
			ctx.AbortWithStatusJSON(http.StatusForbidden,
				helpers.NewHttpResponse(
					http.StatusForbidden,
					"Access Denied - Already Authenticated",
					map[string]interface{}{
						"status":     "logged_in",
						"user_id":    ctx.GetInt("id"),
						"role":       ctx.GetString("role"),
						"message":    "This endpoint is only accessible for non-authenticated users",
						"suggestion": "Please logout first to access this endpoint",
						"metadata": map[string]interface{}{
							"timestamp":      time.Now(),
							"request_path":   ctx.Request.URL.Path,
							"request_method": ctx.Request.Method,
						},
					}))
			return
		}
	}
}

func (m *UserAuthMiddleware) EnsureAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isLogged := common.GetUserStatus(ctx)
		role := ctx.GetString("role")

		if role == strconv.Itoa(int(model.AdminRole)) && isLogged {
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusForbidden, helpers.NewHttpResponse(
				http.StatusForbidden,
				"Admin Access Required",
				map[string]interface{}{
					"message":       "This endpoint requires administrative privileges",
					"current_role":  role,
					"required_role": "admin",
					"metadata": map[string]interface{}{
						"timestamp":      time.Now(),
						"request_path":   ctx.Request.URL.Path,
						"request_method": ctx.Request.Method,
					},
				}))
			return
		}
	}
}

func (m *UserAuthMiddleware) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		is_logged := common.GetUserStatus(ctx)

		if !is_logged {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewHttpResponse(
				http.StatusBadRequest,
				"Logout Failed",
				map[string]interface{}{
					"message":    "No active session found",
					"suggestion": "You must be logged in to perform a logout",
					"metadata": map[string]interface{}{
						"timestamp":      time.Now(),
						"session_status": "inactive",
					},
				}))
			return
		}

		accessToken, err := auth_helper.GetHeader(ctx, auth_helper.AccessTokenHeader)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewHttpResponse(
				http.StatusBadRequest,
				"Invalid Access Token",
				map[string]interface{}{
					"message": "Access token not found in request",
					"error":   err.Error(),
					"metadata": map[string]interface{}{
						"timestamp":      time.Now(),
						"missing_header": auth_helper.AccessTokenHeader,
					},
				}))
			return
		}

		auth_helper.DeleteHeader(ctx, auth_helper.AccessTokenHeader)

		if err := m.AuthManager.SetAccessTokenIntoBlacklist(ctx, accessToken, auth_manager.AccessTokenExpr); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, helpers.NewHttpResponse(
				http.StatusInternalServerError,
				"Access Token Blacklisting Failed",
				map[string]interface{}{
					"message": "Unable to blacklist access token",
					"error":   err.Error(),
					"metadata": map[string]interface{}{
						"timestamp": time.Now(),
					},
				}))
			return
		}

		ctx.Next()
	}
}
