package http

import (
	"net/http"
	"strconv"

	"api/internal/domain"
	"api/internal/errors"
	"api/internal/handler/http/dto"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		errors.MapDomainError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accToken, refToken, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		errors.MapDomainError(c, err)
		return
	}

	c.SetCookie("refresh_token", refToken, 7*24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accToken,
		"message":      "Logged in successfully",
	})
}

func (h *UserHandler) Refresh(c *gin.Context) {
	refreshToken, errCookie := c.Cookie("refresh_token")

	if errCookie != nil || refreshToken == "" {
		var req dto.RefreshRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required via cookie or body"})
			return
		}
		refreshToken = req.RefreshToken
	}

	accToken, refToken, err := h.userService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		errors.MapDomainError(c, err)
		return
	}

	c.SetCookie("refresh_token", refToken, 7*24*60*60, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accToken,
	})
}

func (h *UserHandler) Verify(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.userService.VerifyAccount(c.Request.Context(), id)
	if err != nil {
		errors.MapDomainError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account verified successfully"})
}

func (h *UserHandler) Logout(c *gin.Context) {
	refreshToken, _ := c.Cookie("refresh_token")
	var req dto.LogoutRequest
	if refreshToken == "" {
		if err := c.ShouldBindJSON(&req); err == nil {
			refreshToken = req.RefreshTokenID
		}
	}

	userID := c.GetInt64("userID")
	jti := c.GetString("jti")

	_ = h.userService.Logout(c.Request.Context(), userID, jti, refreshToken)

	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetInt64("userID")

	user, err := h.userService.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
