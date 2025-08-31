package handlers

import (
	"net/http"

	"coderoulette/internal/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	matchService     *services.MatchService
	problemService   *services.ProblemService
	judgeService     *services.JudgeService
	reportService    *services.ReportService
	skillCardService *services.SkillCardService
}

func NewHandlers(
	matchService *services.MatchService,
	problemService *services.ProblemService,
	judgeService *services.JudgeService,
	reportService *services.ReportService,
	skillCardService *services.SkillCardService,
) *Handlers {
	return &Handlers{
		matchService:     matchService,
		problemService:   problemService,
		judgeService:     judgeService,
		reportService:    reportService,
		skillCardService: skillCardService,
	}
}

func (h *Handlers) SetupRoutes(router *gin.Engine) {
	// API routes
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", h.healthCheck)

		// Match routes
		matches := api.Group("/matches")
		{
			matches.POST("/queue", h.queueForMatch)
			matches.GET("/status/:id", h.getMatchStatus)
			matches.GET("/queue-status", h.getQueueStatus)
		}

		// Problem routes
		problems := api.Group("/problems")
		{
			problems.GET("/random", h.getRandomProblem)
			problems.GET("/:id", h.getProblem)
			problems.GET("/", h.getProblems)
			problems.POST("/", h.createProblem)
		}

		// Submission routes
		submissions := api.Group("/submissions")
		{
			submissions.POST("/", h.submitCode)
			submissions.GET("/:id", h.getSubmission)
			submissions.GET("/match/:matchId", h.getMatchSubmissions)
		}

		// Report routes
		reports := api.Group("/reports")
		{
			reports.GET("/:matchId", h.getReport)
			reports.GET("/user/:userId", h.getUserReports)
			reports.GET("/leaderboard", h.getLeaderboard)
		}

		// Skill card routes
		skillCards := api.Group("/skill-cards")
		{
			skillCards.GET("/", h.getAvailableCards)
			skillCards.GET("/player/:playerId", h.getPlayerCards)
			skillCards.POST("/use", h.useSkillCard)
		}
	}

	// WebSocket routes
	router.GET("/ws/match/:roomId", h.handleWebSocket)
}

func (h *Handlers) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "CodeRoulette API is running",
	})
}
