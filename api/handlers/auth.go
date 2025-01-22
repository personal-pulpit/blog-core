package handlers

import (
	"blog/api/helpers"
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/service/authentication"
	"blog/utils"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService authentication.AuthService
}

func NewAuthHandler(authService authentication.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

type (
	verifyEmailInput struct {
		OTP    string `form:"otp" binding:"required"`
		UserID string `form:"userID" binding:"required"`
	}
	
	signinInput struct {
		FirstName string `form:"firstName" binding:"required"`
		LastName  string `form:"lastName" binding:"required"`
		Password  string `form:"password" binding:"required"`
		Email     string `form:"email" binding:"required,emailvalidatior"`
		Biography string `form:"biography" binding:"required"`
	}

	loginInput struct {
		Email    string `form:"email" binding:"required,emailvalidatior"`
		Password string `form:"password" binding:"required"`
	}

	changePasswordInput struct {
		OldPassword string `form:"old_password" binding:"required"`
		NewPassword string `form:"new_password" binding:"required"`
	}
)

func (h *AuthHandler) Register(ctx *gin.Context) {
	var si signinInput
	err := ctx.ShouldBind(&si)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
				http.StatusBadRequest,
				"Missing required registration fields",
				map[string]interface{}{
					"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
					"provided_fields":   si,
				}))
			return
		} else if utils.CheckErrorForWord(err, "emailvalidatior") {
			ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
				http.StatusBadRequest,
				"Invalid email format",
				map[string]interface{}{
					"validation_errors": utils.GetValidationError(ErrInvalidEmail),
					"provided_email":    si.Email,
				}))
			return
		}
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid registration data",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
				"provided_data":     si,
			}))
		return
	}

	user, err := h.AuthService.Register(
		ctx,
		si.FirstName,
		si.LastName,
		si.Email,
		si.Biography,
		si.Password,
	)

	if err != nil {
		if errors.Is(err, postgres_repository.ErrEmailAlreadyExits) ||
			errors.Is(err, postgres_repository.ErrUsernameAlreadyExits) ||
			errors.Is(err, postgres_repository.ErrPhoneNumberAlreadyExits) {
			ctx.JSON(http.StatusConflict, helpers.NewHttpResponse(
				http.StatusConflict,
				"Registration failed - Duplicate information",
				map[string]interface{}{
					"error": err.Error(),
					"email": si.Email,
				}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.NewHttpResponse(
			http.StatusInternalServerError,
			"Registration failed",
			map[string]interface{}{
				"error": err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusCreated, helpers.NewHttpResponse(
		http.StatusCreated,
		"User registered successfully",
		map[string]interface{}{
			"user": map[string]interface{}{
				"id":         user.ID,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"created_at": user.CreatedAt,
			},
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
				"status":    "pending_verification",
			},
		}))
}

func (h *AuthHandler) VerifyEmail(ctx *gin.Context) {
	var input verifyEmailInput
	err := ctx.ShouldBind(&input)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
				http.StatusBadRequest,
				"Missing verification fields",
				map[string]interface{}{
					"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
					"required_fields":   []string{"otp", "userID"},
				}))
			return
		}
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid verification data",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
			}))
		return
	}

	err = h.AuthService.VerifyEmail(ctx, input.OTP, uint(helpers.StringToInt(input.UserID)))
	if err != nil {
		if errors.Is(err, postgres_repository.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, helpers.NewHttpResponse(
				http.StatusNotFound,
				"User not found",
				map[string]interface{}{
					"user_id": input.UserID,
					"error":   err.Error(),
				}))
			return
		}
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Email verification failed",
			map[string]interface{}{
				"error":   err.Error(),
				"user_id": input.UserID,
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Email verified successfully",
		map[string]interface{}{
			"verification_status": "completed",
			"user_id":             input.UserID,
			"metadata": map[string]interface{}{
				"timestamp":   time.Now(),
				"verified_at": time.Now(),
			},
		}))
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var li loginInput
	err := ctx.ShouldBind(&li)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
				http.StatusBadRequest,
				"Missing login credentials",
				map[string]interface{}{
					"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
					"required_fields":   []string{"email", "password"},
				}))
			return
		}
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid login data",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
			}))
		return
	}

	user, accessToken, refreshToken, err := h.AuthService.Login(ctx, li.Email, li.Password)
	if err != nil {
		if errors.Is(err, postgres_repository.ErrUserNotFound) ||
			errors.Is(err, postgres_repository.ErrEmailOrPasswordWrong) {
			ctx.JSON(http.StatusUnauthorized, helpers.NewHttpResponse(
				http.StatusUnauthorized,
				"Invalid credentials",
				map[string]interface{}{
					"error": err.Error(),
				}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.NewHttpResponse(
			http.StatusInternalServerError,
			"Login failed",
			map[string]interface{}{
				"error": err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Login successful",
		map[string]interface{}{
			"user": map[string]interface{}{
				"id":         user.ID,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
			"tokens": map[string]interface{}{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
			"metadata": map[string]interface{}{
				"timestamp":       time.Now(),
				"session_started": time.Now(),
			},
		}))
}

func (h *AuthHandler) Authenticate(ctx *gin.Context) {
	accessToken := ctx.GetHeader("X-Access-Token")

	user, err := h.AuthService.Authenticate(ctx, accessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, helpers.NewHttpResponse(
			http.StatusUnauthorized,
			"Authentication failed",
			map[string]interface{}{
				"error": err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Authentication successful",
		map[string]interface{}{
			"user": map[string]interface{}{
				"id":         user.ID,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
		}))
}

func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	refreshToken := ctx.GetHeader("X-Refresh-Token")
	accessToken := ctx.GetHeader("X-Access-Token")

	newAccessToken, err := h.AuthService.RefreshToken(ctx, refreshToken, accessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, helpers.NewHttpResponse(
			http.StatusUnauthorized,
			"Token refresh failed",
			map[string]interface{}{
				"error": err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Token refreshed successfully",
		map[string]interface{}{
			"access_token": newAccessToken,
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		}))
}

func (h *AuthHandler) ChangePassword(ctx *gin.Context) {
	var input changePasswordInput
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid password change data",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
			}))
		return
	}

	userID := ctx.GetInt("id")

	err := h.AuthService.ChangePassword(ctx, uint(userID), input.OldPassword, input.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Password change failed",
			map[string]interface{}{
				"error": err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Password changed successfully",
		map[string]interface{}{
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		}))
}

func (h *AuthHandler) SendResetPasswordVerification(ctx *gin.Context) {
	type resetPasswordInput struct {
		Email string `form:"email" binding:"required,emailvalidatior"`
	}

	var input resetPasswordInput
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid email format",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
			}))
		return
	}

	token, err := h.AuthService.SendResetPasswordVerification(ctx, input.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Reset password verification failed",
			map[string]interface{}{
				"error": err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Reset password verification sent",
		map[string]interface{}{
			"token": token,
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		}))
}

func (h *AuthHandler) SubmitResetPassword(ctx *gin.Context) {
	type submitResetPasswordInput struct {
		Token       string `form:"token" binding:"required"`
		NewPassword string `form:"new_password" binding:"required"`
	}

	var input submitResetPasswordInput
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid reset password data",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
			}))
		return
	}

	err := h.AuthService.SubmitResetPassword(ctx, input.Token, input.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Reset password failed",
			map[string]interface{}{
				"error": err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Password reset successfully",
		map[string]interface{}{
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		}))
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Logged out successfully",
		map[string]interface{}{
			"metadata": map[string]interface{}{
				"timestamp":     time.Now(),
				"session_ended": time.Now(),
			},
		}))
}
