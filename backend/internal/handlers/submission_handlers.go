package handlers

import (
	"net/http"

	"coderoulette/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubmitCodeRequest struct {
	MatchID   uuid.UUID           `json:"match_id" binding:"required"`
	PlayerID  uuid.UUID           `json:"player_id" binding:"required"`
	Code      string              `json:"code" binding:"required"`
	Language  string              `json:"language" binding:"required"`
	TestCases []services.TestCase `json:"test_cases" binding:"required"`
}

// submitCode handles code submission for judging
func (h *Handlers) submitCode(c *gin.Context) {
	var req SubmitCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	if req.Language == "" {
		req.Language = "go"
	}

	if len(req.TestCases) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "test cases are required"})
		return
	}

	// Submit code for judging
	ctx := c.Request.Context()
	result, err := h.judgeService.SubmitCode(ctx, req.MatchID, req.PlayerID, req.Code, req.Language, req.TestCases)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// getSubmission returns a submission by ID
func (h *Handlers) getSubmission(c *gin.Context) {
	submissionIDStr := c.Param("id")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid submission ID"})
		return
	}

	submission, err := h.judgeService.GetSubmission(submissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "submission not found"})
		return
	}

	c.JSON(http.StatusOK, submission)
}

// getMatchSubmissions returns all submissions for a match
func (h *Handlers) getMatchSubmissions(c *gin.Context) {
	matchIDStr := c.Param("matchId")
	matchID, err := uuid.Parse(matchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match ID"})
		return
	}

	submissions, err := h.judgeService.GetMatchSubmissions(matchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, submissions)
}
