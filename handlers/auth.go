package handlers

import (
	"crypto/rand"
	"encoding/hex"
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

func generateStateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *AuthHandler) Login(c *gin.Context) {
	state := generateStateToken()

	// Optional: Store state in a secure, short-lived HTTP-only cookie to verify in Callback
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	url := h.service.GetAuthURL(state)
	c.JSON(http.StatusOK, gin.H{
		"auth_url": url,
		"state":    state,
	})
}

func (h *AuthHandler) Callback(c *gin.Context) {
	returnedState := c.Query("state")
	savedState, err := c.Cookie("oauth_state")

	if err != nil || returnedState == "" || returnedState != savedState {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid state token (possible CSRF attack)"})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func (h *AuthHandler) SetupGroup(r *gin.RouterGroup) {
	r.GET("/google/login", h.Login)
	r.GET("/google/callback", h.Callback)
}
