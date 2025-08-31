package handlers

import (
	"net/http"

	"coderoulette/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// queueForMatch handles user queueing for matchmaking
func (h *Handlers) queueForMatch(c *gin.Context) {
	var req services.MatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	if req.Difficulty == "" {
		req.Difficulty = "medium"
	}

	if req.Language == "" {
		req.Language = "go"
	}

	// Add user to queue
	ctx := c.Request.Context()
	if err := h.matchService.QueueUser(ctx, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Try to find a match
	result, err := h.matchService.FindMatch(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result == nil {
		// No match found, user is in queue
		c.JSON(http.StatusOK, gin.H{
			"status":  "queued",
			"message": "Waiting for opponent...",
		})
		return
	}

	// Match found
	c.JSON(http.StatusOK, gin.H{
		"status": "matched",
		"match":  result,
	})
}

// getMatchStatus returns the current status of a match
func (h *Handlers) getMatchStatus(c *gin.Context) {
	matchIDStr := c.Param("id")
	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match ID"})
		return
	}

	ctx := c.Request.Context()
	match, err := h.matchService.GetMatchStatus(ctx, matchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	c.JSON(http.StatusOK, match)
}

// getQueueStatus returns the number of users waiting in queue
func (h *Handlers) getQueueStatus(c *gin.Context) {
	difficulty := c.DefaultQuery("difficulty", "medium")
	language := c.DefaultQuery("language", "go")

	ctx := c.Request.Context()
	count, err := h.matchService.GetQueueStatus(ctx, difficulty, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"difficulty": difficulty,
		"language":   language,
		"queue_size": count,
	})
}
