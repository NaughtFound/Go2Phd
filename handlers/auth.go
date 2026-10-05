package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"naughtfound.github.io/go2phd/services"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	url := h.service.GetAuthURL()
	c.JSON(http.StatusOK, gin.H{
		"auth_url": url,
	})
}

func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is missing"})
		return
	}

	token, err := h.service.ExchangeCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token.AccessToken, "refresh_token": token.RefreshToken, "token_type": token.TokenType, "expiry": token.Expiry})
}

func (h *AuthHandler) SetupGroup(r *gin.RouterGroup) {
	r.GET("/google/login", h.Login)
	r.GET("/google/callback", h.Callback)
}
