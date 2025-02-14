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
		OTP    string `json:"otp" binding:"required"`
		UserID string `json:"userID" binding:"required"`
	}

	signinInput struct {
		FirstName string `json:"firstName" binding:"required,min=1,max=50"`
		LastName  string `json:"lastName" binding:"max=25"`
		Password  string `json:"password" binding:"required,min=8,max=200"`
		Email     string `json:"email" binding:"required,emailvalidatior"`
		Biography string `json:"biography" binding:"max=250"`
	}

	loginInput struct {
		Email    string `json:"email" binding:"required,emailvalidatior"`
		Password string `json:"password" binding:"required,min=8,max=200"`
	}

	refreshToken struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	changePasswordInput struct {
		OldPassword string `json:"old_password" binding:"required,min=8,max=200"`
		NewPassword string `json:"new_password" binding:"required,min=8,max=200"`
	}

	resetPasswordInput struct {
		Email string `json:"email" binding:"required,emailvalidatior"`
	}

	submitResetPasswordInput struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	logoutInput struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
)

// @Summary Register a new user
// @Description Create a new user account with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param user body signinInput true "User registration data"
// @Success 201 {object} helpers.HttpResponse{data=map[string]interface{}} "User registered successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid registration data"
// @Failure 409 {object} helpers.HttpResponse{data=map[string]interface{}} "Email already exists"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var si signinInput
	err := ctx.ShouldBindJSON(&si)
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
		if errors.Is(err, postgres_repository.ErrEmailAlreadyExits) {
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

// @Summary Verify user email
// @Description Verify the email address of a user using the provided OTP
// @Tags auth
// @Accept json
// @Produce json
// @Param input body verifyEmailInput true "Email verification data"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Email verified successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid verification data"
// @Failure 404 {object} helpers.HttpResponse{data=map[string]interface{}} "User not found"
// @Router /api/v1/auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(ctx *gin.Context) {
	var input verifyEmailInput
	err := ctx.ShouldBindJSON(&input)
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

// @Summary User login
// @Description Authenticate a user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body loginInput true "Login credentials"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Login successful"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid login data"
// @Failure 401 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid credentials"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var li loginInput
	err := ctx.ShouldBindJSON(&li)
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

	_, accessToken, refreshToken, err := h.AuthService.Login(ctx, li.Email, li.Password)
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

// @Summary Authenticate user
// @Description Validate the access token of the user
// @Tags auth
// @Produce json
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Authentication successful"
// @Failure 401 {object} helpers.HttpResponse{data=map[string]interface{}} "Authentication failed"
// @Router /api/v1/auth/authenticate [get]
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

// @Summary Refresh access token
// @Description Refresh the user's access token using the refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body refreshToken true "Refresh token data"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Token refreshed successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid request data"
// @Failure 401 {object} helpers.HttpResponse{data=map[string]interface{}} "Token refresh failed"
// @Router /api/v1/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var li = new(refreshToken)

	err := ctx.ShouldBindJSON(&li)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
				http.StatusBadRequest,
				"Missing refresh token",
				map[string]interface{}{
					"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
					"required_fields":   []string{"refresh_token"},
				}))
			return
		}
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid request data",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
			}))
		return
	}
	
	accessToken := ctx.GetString("accessToken")

	newAccessToken, err := h.AuthService.RefreshToken(ctx, li.RefreshToken, accessToken)
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
			"new_access_token": newAccessToken,
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		}))
}

// @Summary Change user password
// @Description Change the password of the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body changePasswordInput true "Password change data"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Password changed successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid password change data"
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(ctx *gin.Context) {
	var input changePasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
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

// @Summary Send reset password verification
// @Description Send a verification email to reset the password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body resetPasswordInput true "Email data"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Reset password verification sent"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid email format"
// @Router /api/v1/auth/reset-password/request [post]
func (h *AuthHandler) SendResetPasswordVerification(ctx *gin.Context) {
	var input resetPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
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

// @Summary Submit reset password
// @Description Submit a new password to reset the user's password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body submitResetPasswordInput true "Reset password data"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Password reset successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid reset password data"
// @Router /api/v1/auth/reset-password/submit [post]
func (h *AuthHandler) SubmitResetPassword(ctx *gin.Context) {
	var input submitResetPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
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

// @Summary User logout
// @Description Log out the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body logoutInput true "Logout data"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Logged out successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid logout data"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	var input logoutInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid logout data",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
			}))
		return
	}

	err := h.AuthService.DestroyRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusInternalServerError,
			"Logout failed:could not destroy refresh token",
			map[string]interface{}{
				"error": err.Error(),
			}))
		return
	}

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
