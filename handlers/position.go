package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"naughtfound.github.io/go2phd/services"
)

type PositionHandler struct {
	positionService *services.PositionService
	statsService    *services.StatsService
}

func NewPositionHandler(positionService *services.PositionService, statsService *services.StatsService) *PositionHandler {
	return &PositionHandler{
		positionService: positionService,
		statsService:    statsService,
	}
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

	positions, err := h.positionService.FetchAllPositions(refreshToken, spreadsheetID, sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, positions)
}

func (h *PositionHandler) GetStats(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	spreadsheetID := c.GetHeader("X-Spreadsheet-ID")
	sheetName := c.DefaultQuery("sheet", "Sheet1")
	groupBy := c.Query("group_by")

	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required header: X-Refresh-Token"})
		return
	}

	if spreadsheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required header: X-Spreadsheet-ID"})
		return
	}

	analytics, err := h.statsService.GetAnalytics(refreshToken, spreadsheetID, sheetName, groupBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func (h *PositionHandler) SetupGroup(r *gin.RouterGroup) {
	r.GET("", h.GetAllPositions)
	r.GET("/stats", h.GetStats)
}
