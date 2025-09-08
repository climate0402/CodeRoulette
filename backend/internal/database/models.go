package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	BaseIDModel
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Rating    int       `gorm:"default:1200" json:"rating"`
	Wins      int       `gorm:"default:0" json:"wins"`
	Losses    int       `gorm:"default:0" json:"losses"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Problem represents a coding problem
type Problem struct {
	BaseIDModel
	Title       string    `gorm:"not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Difficulty  string    `gorm:"not null" json:"difficulty"`   // easy, medium, hard
	Language    string    `gorm:"not null" json:"language"`     // go, python, javascript
	TestCases   string    `gorm:"type:jsonb" json:"test_cases"` // JSON array of test cases
	Solution    string    `gorm:"type:text" json:"solution"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Match represents a match between two players
type Match struct {
	BaseIDModel
	Player1ID uuid.UUID  `gorm:"not null" json:"player1_id"`
	Player2ID uuid.UUID  `gorm:"not null" json:"player2_id"`
	ProblemID uuid.UUID  `gorm:"not null" json:"problem_id"`
	Status    string     `gorm:"default:'waiting'" json:"status"` // waiting, active, completed
	WinnerID  *uuid.UUID `json:"winner_id"`
	Duration  int        `json:"duration"` // in seconds
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relations
	Player1 User    `gorm:"foreignKey:Player1ID" json:"player1"`
	Player2 User    `gorm:"foreignKey:Player2ID" json:"player2"`
	Problem Problem `gorm:"foreignKey:ProblemID" json:"problem"`
	Winner  *User   `gorm:"foreignKey:WinnerID" json:"winner"`
}

// Submission represents a code submission by a player 记录对战过程，提供“回放”给观众，同时report总结时可以查看提交的所有信息
type Submission struct {
	BaseIDModel
	MatchID   uuid.UUID `gorm:"not null" json:"match_id"`
	PlayerID  uuid.UUID `gorm:"not null" json:"player_id"`
	Code      string    `gorm:"type:text;not null" json:"code"`
	Language  string    `gorm:"not null" json:"language"`
	Status    string    `gorm:"default:'pending'" json:"status"` // pending, running, passed, failed
	Score     int       `json:"score"`
	Runtime   int       `json:"runtime"` // in milliseconds
	ErrorMsg  string    `gorm:"type:text" json:"error_msg"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Match  Match `gorm:"foreignKey:MatchID" json:"match"`
	Player User  `gorm:"foreignKey:PlayerID" json:"player"`
}

// Report represents a match report
type Report struct {
	BaseIDModel
	MatchID   uuid.UUID `gorm:"not null" json:"match_id"`
	Data      string    `gorm:"type:jsonb;not null" json:"data"` // JSON report data
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Match Match `gorm:"foreignKey:MatchID" json:"match"`
}

// SkillCard represents a skill card that can be used in matches
type SkillCard struct {
	BaseIDModel
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Type        string    `gorm:"not null" json:"type"` // peek, swap, hint, etc.
	Cost        int       `gorm:"default:1" json:"cost"`
	Rarity      string    `gorm:"default:'common'" json:"rarity"` // common, rare, epic, legendary
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate hook for UUID generation
//统一写法
type BaseIDModel struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
}

func (bm *BaseIDModel) BeforeCreate(tx *gorm.DB) error {
	if bm.ID == uuid.Nil {
		bm.ID = uuid.New()
	}
	return nil
}