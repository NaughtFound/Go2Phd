package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"naughtfound.github.io/go2phd/services"
)

type UniversityHandler struct {
	service *services.UniversityService
}

func NewUniversityHandler(service *services.UniversityService) *UniversityHandler {
	return &UniversityHandler{service: service}
}

func (h *UniversityHandler) GetUniversityInfo(c *gin.Context) {
	name, _ := c.Params.Get("name")
	user, err := h.service.GetUniversityInfo(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UniversityHandler) SetupGroup(r *gin.RouterGroup) {
	r.GET("/:name", h.GetUniversityInfo)
}
