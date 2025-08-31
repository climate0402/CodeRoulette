package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UseSkillCardRequest struct {
	MatchID  string `json:"match_id" binding:"required"`
	PlayerID string `json:"player_id" binding:"required"`
	CardID   string `json:"card_id" binding:"required"`
}

// getAvailableCards returns all available skill cards
func (h *Handlers) getAvailableCards(c *gin.Context) {
	cards := h.skillCardService.GetAvailableCards()
	c.JSON(http.StatusOK, gin.H{
		"cards": cards,
	})
}

// getPlayerCards returns skill cards available to a player
func (h *Handlers) getPlayerCards(c *gin.Context) {
	playerID := c.Param("playerId")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player ID is required"})
		return
	}

	ctx := c.Request.Context()
	cards, err := h.skillCardService.GetPlayerCards(ctx, playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"player_id": playerID,
		"cards":     cards,
	})
}

// useSkillCard uses a skill card in a match
func (h *Handlers) useSkillCard(c *gin.Context) {
	var req UseSkillCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	usage, err := h.skillCardService.UseSkillCard(ctx, req.MatchID, req.PlayerID, req.CardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Skill card used successfully",
		"usage":   usage,
	})
}
