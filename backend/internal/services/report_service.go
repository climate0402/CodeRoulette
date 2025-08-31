package services

import (
	"encoding/json"
	"time"

	"coderoulette/internal/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportService struct {
	db *gorm.DB
}

type ReportData struct {
	MatchID     uuid.UUID           `json:"match_id"`
	Player1     PlayerStats         `json:"player1"`
	Player2     PlayerStats         `json:"player2"`
	Winner      *uuid.UUID          `json:"winner"`
	Duration    int                 `json:"duration"`
	Problem     ProblemSummary      `json:"problem"`
	Submissions []SubmissionSummary `json:"submissions"`
	CreatedAt   time.Time           `json:"created_at"`
}

type PlayerStats struct {
	ID               uuid.UUID `json:"id"`
	Username         string    `json:"username"`
	FinalScore       int       `json:"final_score"`
	BestScore        int       `json:"best_score"`
	TotalSubmissions int       `json:"total_submissions"`
	AverageRuntime   int       `json:"average_runtime"`
	FirstSubmission  time.Time `json:"first_submission"`
	LastSubmission   time.Time `json:"last_submission"`
}

type ProblemSummary struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Difficulty    string    `json:"difficulty"`
	Language      string    `json:"language"`
	TestCaseCount int       `json:"test_case_count"`
}

type SubmissionSummary struct {
	ID        uuid.UUID `json:"id"`
	PlayerID  uuid.UUID `json:"player_id"`
	Score     int       `json:"score"`
	Runtime   int       `json:"runtime"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

// GenerateReport generates a comprehensive report for a match
func (s *ReportService) GenerateReport(matchID uuid.UUID) (*ReportData, error) {
	// Get match details
	var match database.Match
	if err := s.db.Preload("Player1").Preload("Player2").Preload("Problem").First(&match, "id = ?", matchID).Error; err != nil {
		return nil, err
	}

	// Get all submissions for this match
	var submissions []database.Submission
	if err := s.db.Where("match_id = ?", matchID).Order("created_at ASC").Find(&submissions).Error; err != nil {
		return nil, err
	}

	// Calculate player statistics
	player1Stats := s.calculatePlayerStats(match.Player1ID, submissions)
	player2Stats := s.calculatePlayerStats(match.Player2ID, submissions)

	// Create submission summaries
	submissionSummaries := make([]SubmissionSummary, len(submissions))
	for i, sub := range submissions {
		submissionSummaries[i] = SubmissionSummary{
			ID:        sub.ID,
			PlayerID:  sub.PlayerID,
			Score:     sub.Score,
			Runtime:   sub.Runtime,
			Status:    sub.Status,
			CreatedAt: sub.CreatedAt,
		}
	}

	// Parse test cases to get count
	var testCases []TestCase
	json.Unmarshal([]byte(match.Problem.TestCases), &testCases)

	report := &ReportData{
		MatchID:  matchID,
		Player1:  player1Stats,
		Player2:  player2Stats,
		Winner:   match.WinnerID,
		Duration: match.Duration,
		Problem: ProblemSummary{
			ID:            match.Problem.ID,
			Title:         match.Problem.Title,
			Difficulty:    match.Problem.Difficulty,
			Language:      match.Problem.Language,
			TestCaseCount: len(testCases),
		},
		Submissions: submissionSummaries,
		CreatedAt:   match.CreatedAt,
	}

	// Save report to database
	if err := s.saveReport(matchID, report); err != nil {
		return nil, err
	}

	return report, nil
}

// calculatePlayerStats calculates statistics for a player
func (s *ReportService) calculatePlayerStats(playerID uuid.UUID, submissions []database.Submission) PlayerStats {
	var playerStats PlayerStats
	playerStats.ID = playerID

	var playerSubmissions []database.Submission
	for _, sub := range submissions {
		if sub.PlayerID == playerID {
			playerSubmissions = append(playerSubmissions, sub)
		}
	}

	if len(playerSubmissions) == 0 {
		return playerStats
	}

	// Get player username
	var player database.User
	s.db.First(&player, "id = ?", playerID)
	playerStats.Username = player.Username

	// Calculate statistics
	playerStats.TotalSubmissions = len(playerSubmissions)
	playerStats.FirstSubmission = playerSubmissions[0].CreatedAt
	playerStats.LastSubmission = playerSubmissions[len(playerSubmissions)-1].CreatedAt

	bestScore := 0
	totalRuntime := 0
	validSubmissions := 0

	for _, sub := range playerSubmissions {
		if sub.Score > bestScore {
			bestScore = sub.Score
		}
		if sub.Runtime > 0 {
			totalRuntime += sub.Runtime
			validSubmissions++
		}
	}

	playerStats.BestScore = bestScore
	playerStats.FinalScore = playerSubmissions[len(playerSubmissions)-1].Score

	if validSubmissions > 0 {
		playerStats.AverageRuntime = totalRuntime / validSubmissions
	}

	return playerStats
}

// saveReport saves the report to database
func (s *ReportService) saveReport(matchID uuid.UUID, report *ReportData) error {
	reportJSON, err := json.Marshal(report)
	if err != nil {
		return err
	}

	dbReport := &database.Report{
		ID:      uuid.New(),
		MatchID: matchID,
		Data:    string(reportJSON),
	}

	return s.db.Create(dbReport).Error
}

// GetReport retrieves a report by match ID
func (s *ReportService) GetReport(matchID uuid.UUID) (*ReportData, error) {
	var report database.Report
	if err := s.db.First(&report, "match_id = ?", matchID).Error; err != nil {
		return nil, err
	}

	var reportData ReportData
	if err := json.Unmarshal([]byte(report.Data), &reportData); err != nil {
		return nil, err
	}

	return &reportData, nil
}

// GetUserReports retrieves all reports for a user
func (s *ReportService) GetUserReports(userID uuid.UUID, page, limit int) ([]*ReportData, int64, error) {
	var reports []database.Report
	var count int64

	// Count total reports for user
	if err := s.db.Table("reports").
		Joins("JOIN matches ON reports.match_id = matches.id").
		Where("matches.player1_id = ? OR matches.player2_id = ?", userID, userID).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get reports with pagination
	offset := (page - 1) * limit
	if err := s.db.Table("reports").
		Joins("JOIN matches ON reports.match_id = matches.id").
		Where("matches.player1_id = ? OR matches.player2_id = ?", userID, userID).
		Offset(offset).Limit(limit).
		Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	// Parse report data
	result := make([]*ReportData, len(reports))
	for i, report := range reports {
		var reportData ReportData
		json.Unmarshal([]byte(report.Data), &reportData)
		result[i] = &reportData
	}

	return result, count, nil
}

// GetLeaderboard returns the top players
func (s *ReportService) GetLeaderboard(limit int) ([]PlayerStats, error) {
	var users []database.User
	if err := s.db.Order("rating DESC").Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	result := make([]PlayerStats, len(users))
	for i, user := range users {
		result[i] = PlayerStats{
			ID:         user.ID,
			Username:   user.Username,
			FinalScore: user.Rating, // Using rating as final score for leaderboard
		}
	}

	return result, nil
}
