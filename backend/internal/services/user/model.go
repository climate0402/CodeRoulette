package user

import (
	"time"

	models "coderoulette/internal/database"
)

type User struct {
	models.BaseIDModel
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Rating    int       `json:"rating"`
	Wins      int       `json:"wins"`
	Losses    int       `json:"losses"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
