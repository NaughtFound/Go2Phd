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
	refreshToken := c.GetHeader("X-Refresh-Token")
	spreadsheetID := c.GetHeader("X-Spreadsheet-ID")
	sheetName := c.DefaultQuery("sheet", "Sheet1")

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required header: X-Refresh-Token"})
		return
	}

	if spreadsheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required header: X-Spreadsheet-ID"})
		return
	}

	positions, err := h.service.FetchAllPositions(refreshToken, spreadsheetID, sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, positions)
}

func (h *PositionHandler) SetupGroup(r *gin.RouterGroup) {
	r.GET("", h.GetAllPositions)
}
