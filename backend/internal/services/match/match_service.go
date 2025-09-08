package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"coderoulette/internal/database"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MatchService struct {
	redis *redis.Client
	db    *gorm.DB
}

type MatchRequest struct {
	UserID     uuid.UUID `json:"user_id"`
	Difficulty string    `json:"difficulty"`
	Language   string    `json:"language"`
}

type MatchResult struct {
	MatchID   uuid.UUID `json:"match_id"`
	Player1ID uuid.UUID `json:"player1_id"`
	Player2ID uuid.UUID `json:"player2_id"`
	ProblemID uuid.UUID `json:"problem_id"`
	RoomID    string    `json:"room_id"`
}

func NewMatchService(redis *redis.Client) *MatchService {
	return &MatchService{
		redis: redis,
	}
}

func (s *MatchService) SetDB(db *gorm.DB) {
	s.db = db
}

// QueueUser adds a user to the matchmaking queue
func (s *MatchService) QueueUser(ctx context.Context, req *MatchRequest) error {
	queueKey := fmt.Sprintf("queue:%s:%s", req.Difficulty, req.Language)

	// Add user to queue with timestamp
	queueData := map[string]interface{}{
		"user_id":    req.UserID.String(),
		"timestamp":  time.Now().Unix(),
		"difficulty": req.Difficulty,
		"language":   req.Language,
	}

	data, err := json.Marshal(queueData)
	if err != nil {
		return err
	}

	// Add to sorted set with timestamp as score for FIFO ordering
	return s.redis.ZAdd(ctx, queueKey, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: data,
	}).Err()
}

// FindMatch attempts to find a match for the user
func (s *MatchService) FindMatch(ctx context.Context, req *MatchRequest) (*MatchResult, error) {
	queueKey := fmt.Sprintf("queue:%s:%s", req.Difficulty, req.Language)

	// Get all users in queue
	users, err := s.redis.ZRange(ctx, queueKey, 0, 10).Result()
	if err != nil {
		return nil, err
	}

	if len(users) < 2 {
		return nil, nil // No match found yet
	}

	// Parse first two users
	var user1, user2 map[string]interface{}
	if err := json.Unmarshal([]byte(users[0]), &user1); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(users[1]), &user2); err != nil {
		return nil, err
	}

	// Remove matched users from queue
	s.redis.ZRem(ctx, queueKey, users[0], users[1])

	// Create match
	matchID := uuid.New()
	player1ID, _ := uuid.Parse(user1["user_id"].(string))
	player2ID, _ := uuid.Parse(user2["user_id"].(string))

	// Create match record in database
	match := &database.Match{
		ID:        matchID,
		Player1ID: player1ID,
		Player2ID: player2ID,
		Status:    "waiting",
	}

	if err := s.db.Create(match).Error; err != nil {
		return nil, err
	}

	// Generate room ID for WebSocket
	roomID := fmt.Sprintf("room:%s", matchID.String())

	return &MatchResult{
		MatchID:   matchID,
		Player1ID: player1ID,
		Player2ID: player2ID,
		RoomID:    roomID,
	}, nil
}

// GetMatchStatus returns the current status of a match
func (s *MatchService) GetMatchStatus(ctx context.Context, matchID uuid.UUID) (*database.Match, error) {
	var match database.Match
	if err := s.db.Preload("Player1").Preload("Player2").Preload("Problem").First(&match, "id = ?", matchID).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

// UpdateMatchStatus updates the status of a match
func (s *MatchService) UpdateMatchStatus(ctx context.Context, matchID uuid.UUID, status string) error {
	return s.db.Model(&database.Match{}).Where("id = ?", matchID).Update("status", status).Error
}

// CompleteMatch marks a match as completed and sets the winner
func (s *MatchService) CompleteMatch(ctx context.Context, matchID uuid.UUID, winnerID *uuid.UUID, duration int) error {
	updates := map[string]interface{}{
		"status":   "completed",
		"duration": duration,
	}

	if winnerID != nil {
		updates["winner_id"] = winnerID
	}

	return s.db.Model(&database.Match{}).Where("id = ?", matchID).Updates(updates).Error
}

// GetQueueStatus returns the number of users waiting in queue
func (s *MatchService) GetQueueStatus(ctx context.Context, difficulty, language string) (int, error) {
	queueKey := fmt.Sprintf("queue:%s:%s", difficulty, language)
	count, err := s.redis.ZCard(ctx, queueKey).Result()
	return int(count), err
}

// RemoveFromQueue removes a user from the matchmaking queue
func (s *MatchService) RemoveFromQueue(ctx context.Context, userID uuid.UUID, difficulty, language string) error {
	queueKey := fmt.Sprintf("queue:%s:%s", difficulty, language)

	// Get all users in queue and find the one to remove
	users, err := s.redis.ZRange(ctx, queueKey, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, userData := range users {
		var user map[string]interface{}
		if err := json.Unmarshal([]byte(userData), &user); err != nil {
			continue
		}

		if user["user_id"].(string) == userID.String() {
			return s.redis.ZRem(ctx, queueKey, userData).Err()
		}
	}

	return nil
}
