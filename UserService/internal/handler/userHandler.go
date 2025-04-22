package handler

import (
	"CarStore/UserService/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	uc *usecase.UserUsecase
}

func NewAuthHandler(rg *gin.RouterGroup, uc *usecase.UserUsecase) {
	h := &AuthHandler{
		uc: uc,
	}
	rg.POST("/register", h.Register)
	rg.POST("/login", h.Login)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required,alphanum"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required,oneof=admin user"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.uc.Register(c.Request.Context(), req.Email, req.Username, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.uc.Login(c.Request.Context(), req.Identifier, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
