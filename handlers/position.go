package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"naughtfound.github.io/go2phd/services"
)

type PositionHandler struct {
	service *services.PositionService
}

func NewPositionHandler(service *services.PositionService) *PositionHandler {
	return &PositionHandler{service: service}
}

func (h *PositionHandler) GetAllPositions(c *gin.Context) {
	sheetName := c.DefaultQuery("sheet", "Sheet1")

	positions, err := h.service.FetchAllPositions(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, positions)
}

func (h *PositionHandler) SetupGroup(r *gin.RouterGroup) {
	r.GET("", h.GetAllPositions)
}
