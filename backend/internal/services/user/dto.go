package user

type UserCreate struct {
    Username string `json:"username binding:"required"`
    Email string `json:"email binding:"required"`
    Password string `json:"password binding:"required"`
}

type UserResponse struct
 {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Rating   int    `json:"rating"`
    Wins     int    `json:"wins"`
    Losses   int    `json:"losses"`
}

type UserLogin struct {
    Username string `json:"username binding:"required"`
    Password string `json:"password binding:"required"`
}