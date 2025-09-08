package user

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type UserHandler struct{
    service *UserService
}

func NewUserHandler(service *UserService)*UserHandler {
    return &UserHandler{service: service}
}

// 注册接口
func (h *UserHandler) Register(c *gin.Context) {
    var req UserCreate
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.service.Register(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    resp := UserResponse{
        ID:       user.ID.String(),
        Username: user.Username,
        Email:    user.Email,
        Rating:   user.Rating,
        Wins:     user.Wins,
        Losses:   user.Losses,
    }

    c.JSON(http.StatusOK, resp)
}

// 登录接口
func (h *UserHandler) Login(c *gin.Context) {
    var req UserLogin
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.service.Login(req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    resp := UserResponse{
        ID:       user.ID.String(),
        Username: user.Username,
        Email:    user.Email,
        Rating:   user.Rating,
        Wins:     user.Wins,
        Losses:   user.Losses,
    }

    c.JSON(http.StatusOK, resp)
}