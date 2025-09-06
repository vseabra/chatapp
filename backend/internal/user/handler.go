package user

import (
	"chatapp/internal/config"
	"chatapp/internal/constants"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1/auth")
	api.POST("/register", h.handleRegister)
	api.POST("/login", h.handleLogin)
}

func (h *Handler) handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.TIMEOUT_SECONDS)
	defer cancel()
	resp, err := h.service.Register(ctx, req)
	if err != nil {
		status := http.StatusInternalServerError
		reason := err.Error()
		if err == ErrEmailTaken {
			status = http.StatusConflict
			reason = "email already in use"
		}
		c.JSON(status, gin.H{"error": reason})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.TIMEOUT_SECONDS)
	defer cancel()
	cfg := config.Load()
	resp, err := h.service.Login(ctx, req, cfg.JWTSecret, cfg.JWTExpires)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
