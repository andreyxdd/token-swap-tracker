package handlers

import (
	"consumer/internal/rest/middleware"
	"consumer/internal/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type StatsHandler struct {
	service *services.StatsService
}

func NewStatsHandler(s *services.StatsService) *StatsHandler {
	return &StatsHandler{s}
}

func (h *StatsHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/stats/:window/tokens/:token", h.GetTokenStats)
	r.GET("/stats/:window/pairs/:pair", h.GetPairStats)
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// @Summary Get single token stats in a specific period window
// @Description Available period winows are "5min", "1h", "24h". Available tokens are "BTC", "USDT", "TON", "SOL", "ETH".
// @Tags Stats
// @Accept json
// @Produce json
// @Success 200 {array} models.Stats "Token stats"
// @Success 400 {array} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stats/{window}/tokens/{token} [get]
func (h *StatsHandler) GetTokenStats(c *gin.Context) {
	ctx := c.Request.Context()

	window := c.Param("window")
	isValidPeriod := middleware.IsValidPeriod(window)
	if !isValidPeriod {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period window provided"})
		return
	}

	token := c.Param("token")
	isValidToken := middleware.IsValidToken(token)
	if !isValidToken {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token provided"})
		return
	}

	statsKey := fmt.Sprintf("stats:%s:%s", token, window)
	stats, err := h.service.GetStats(ctx, statsKey)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to get stats from the stats service for the key %s", statsKey).Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if stats == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no stats for the provided token and period window"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary Get swap pair stats in a specific period window
// @Description Available period are "5min", "1h", "24h". Available tokens are "BTC", "USDT", "TON", "SOL", "ETH".
// @Tags Stats
// @Accept json
// @Produce json
// @Success 200 {array} models.Stats "Swap pair stats"
// @Success 400 {array} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /stats/{window}/pairs/{pair} [get]
func (h *StatsHandler) GetPairStats(c *gin.Context) {
	ctx := c.Request.Context()

	window := c.Param("window")
	isValidPeriod := middleware.IsValidPeriod(window)
	if !isValidPeriod {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period window provided"})
		return
	}

	pair := c.Param("pair")
	isValidPair := middleware.IsValidPair(pair)
	if !isValidPair {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pair provided"})
		return
	}

	statsKey := fmt.Sprintf("stats:%s:%s", pair, window)
	stats, err := h.service.GetStats(ctx, statsKey)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "failed to get stats from the stats service for the key %s", statsKey).Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	if stats == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no stats for the provided pair and period window"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
