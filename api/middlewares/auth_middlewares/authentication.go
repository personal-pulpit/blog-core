package auth_middlewares

import (
	"blog/api/helpers"
	"blog/api/helpers/auth_helper"
	"blog/api/helpers/common"
	"blog/internal/model"
	"blog/pkg/auth_manager"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthManager interface {
	IsAccessTokenBlacklisted(ctx context.Context, accessToken string) bool
	DecodeAccessToken(ctx context.Context, accessToken string) (*auth_manager.AccessTokenClaims, error)
	SetAccessTokenIntoBlacklist(ctx context.Context, accessToken string, expiresAt time.Duration) error
}

type UserAuthMiddleware struct {
	AuthManager AuthManager
}

func NewUserAuthMiddleware(authManger AuthManager) *UserAuthMiddleware {
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
				helpers.RespondWithError(ctx, http.StatusUnauthorized, "Token Blacklisted", map[string]interface{}{
					"message":    "Your session token has been invalidated",
					"suggestion": "Please login again to obtain a new token",
					"metadata": map[string]interface{}{
						"timestamp":    time.Now(),
						"token_status": "blacklisted",
					},
				})
				return
			}

			accessTokenClaims, err := m.AuthManager.DecodeAccessToken(ctx, accessToken)
			if err != nil {
				if ctx.Request.URL.Path == "/api/v1/auth/refresh-token" {
					ctx.Set("is_logged", false)
					ctx.Set("accessToken", accessToken)
					return
				}

				helpers.RespondWithError(ctx, http.StatusUnauthorized, "Invalid Token", map[string]interface{}{
					"message":    "Unable to verify your authentication token",
					"suggestion": "Please ensure your token is valid or login again",
					"error":      err.Error(),
					"metadata": map[string]interface{}{
						"timestamp":    time.Now(),
						"token_status": "invalid",
					},
				})
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
			helpers.RespondWithError(ctx, http.StatusUnauthorized, "Authentication Required", map[string]interface{}{
				"message":    "You must be logged in to access this resource",
				"suggestion": "Please login to continue",
				"metadata": map[string]interface{}{
					"timestamp":      time.Now(),
					"request_path":   ctx.Request.URL.Path,
					"request_method": ctx.Request.Method,
				},
			})
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
			helpers.RespondWithError(ctx, http.StatusForbidden, "Access Denied - Already Authenticated", map[string]interface{}{
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
			})
			return
		}
	}
}

func (m *UserAuthMiddleware) EnsureAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isLogged := common.GetUserStatus(ctx)
		role := ctx.GetString("role")

		if common.IsAdmin(model.Role(helpers.StringToInt(role))) && isLogged {
			ctx.Next()
		} else {
			helpers.RespondWithError(ctx, http.StatusForbidden, "Admin Access Required", map[string]interface{}{
				"message":       "This endpoint requires administrative privileges",
				"current_role":  role,
				"required_role": "admin",
				"metadata": map[string]interface{}{
					"timestamp":      time.Now(),
					"request_path":   ctx.Request.URL.Path,
					"request_method": ctx.Request.Method,
				},
			})
			return
		}
	}
}

func (m *UserAuthMiddleware) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		is_logged := common.GetUserStatus(ctx)

		if !is_logged {
			helpers.RespondWithError(ctx, http.StatusBadRequest, "Logout Failed", map[string]interface{}{
				"message":    "No active session found",
				"suggestion": "You must be logged in to perform a logout",
				"metadata": map[string]interface{}{
					"timestamp":      time.Now(),
					"session_status": "inactive",
				},
			})
			return
		}

		accessToken, err := auth_helper.GetHeader(ctx, auth_helper.AccessTokenHeader)
		if err != nil {
			helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid Access Token", map[string]interface{}{
				"message": "Access token not found in request",
				"error":   err.Error(),
				"metadata": map[string]interface{}{
					"timestamp":      time.Now(),
					"missing_header": auth_helper.AccessTokenHeader,
				},
			})
			return
		}

		auth_helper.DeleteHeader(ctx, auth_helper.AccessTokenHeader)

		if err := m.AuthManager.SetAccessTokenIntoBlacklist(ctx, accessToken, auth_manager.AccessTokenExpr); err != nil {
			helpers.RespondWithError(ctx, http.StatusInternalServerError, "Access Token Blacklisting Failed", map[string]interface{}{
				"message": "Unable to blacklist access token",
				"error":   err.Error(),
				"metadata": map[string]interface{}{
					"timestamp": time.Now(),
				},
			})
			return
		}

		ctx.Next()
	}
}
