package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
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
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
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
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
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

// Submission represents a code submission by a player
type Submission struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
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
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	MatchID   uuid.UUID `gorm:"not null" json:"match_id"`
	Data      string    `gorm:"type:jsonb;not null" json:"data"` // JSON report data
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Match Match `gorm:"foreignKey:MatchID" json:"match"`
}

// SkillCard represents a skill card that can be used in matches
type SkillCard struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Type        string    `gorm:"not null" json:"type"` // peek, swap, hint, etc.
	Cost        int       `gorm:"default:1" json:"cost"`
	Rarity      string    `gorm:"default:'common'" json:"rarity"` // common, rare, epic, legendary
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate hook for UUID generation
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (p *Problem) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (m *Match) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (s *Submission) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

func (r *Report) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (sc *SkillCard) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	return nil
}
