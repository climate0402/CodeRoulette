package handlers

import (
	"net/http"
	"strconv"

	"coderoulette/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getRandomProblem returns a random problem
func (h *Handlers) getRandomProblem(c *gin.Context) {
	difficulty := c.DefaultQuery("difficulty", "medium")
	language := c.DefaultQuery("language", "go")

	problem, err := h.problemService.GetRandomProblem(difficulty, language)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No problems found"})
		return
	}

	c.JSON(http.StatusOK, problem)
}

// getProblem returns a specific problem by ID
func (h *Handlers) getProblem(c *gin.Context) {
	problemIDStr := c.Param("id")
	problemID, err := uuid.Parse(problemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid problem ID"})
		return
	}

	problem, err := h.problemService.GetProblemByID(problemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "problem not found"})
		return
	}

	c.JSON(http.StatusOK, problem)
}

// getProblems returns a list of problems with pagination
func (h *Handlers) getProblems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	difficulty := c.Query("difficulty")
	language := c.Query("language")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	problems, total, err := h.problemService.GetProblems(page, limit, difficulty, language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"problems": problems,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// createProblem creates a new problem
func (h *Handlers) createProblem(c *gin.Context) {
	var problem services.ProblemData
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if problem.Title == "" || problem.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and description are required"})
		return
	}

	if problem.Difficulty == "" {
		problem.Difficulty = "medium"
	}

	if problem.Language == "" {
		problem.Language = "go"
	}

	if err := h.problemService.CreateProblem(&problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, problem)
}
