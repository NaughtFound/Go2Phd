package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"naughtfound.github.io/go2phd/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	user, err := h.service.GetUser(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) SetupGroup(r *gin.RouterGroup) {
	r.GET("/:id", h.GetUser)
}
