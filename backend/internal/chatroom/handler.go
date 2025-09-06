package chatroom

import (
	"chatapp/internal/auth"
	"chatapp/internal/constants"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine, jwtSecret string) {
	group := r.Group(constants.APIv1 + "/chatroom")
	group.Use(auth.AuthMiddleware(jwtSecret))
	group.POST("", h.create)
	group.GET("all", h.listAll)
	group.PUT(":id", h.rename)
	group.DELETE(":id", h.delete)
}

func (h *Handler) create(c *gin.Context) {
	var req CreateChatRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	uid := c.GetString("uid")
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.TIMEOUT_SECONDS)
	defer cancel()
	resp, err := h.service.Create(ctx, uid, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) listAll(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	skip, _ := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.TIMEOUT_SECONDS)
	defer cancel()
	items, err := h.service.ListAll(ctx, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) rename(c *gin.Context) {
	uid := c.GetString("uid")
	id := c.Param("id")
	var req UpdateChatRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.TIMEOUT_SECONDS)
	defer cancel()
	if err := h.service.Rename(ctx, uid, id, req); err != nil {
		status := http.StatusBadRequest
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) delete(c *gin.Context) {
	uid := c.GetString("uid")
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.TIMEOUT_SECONDS)
	defer cancel()
	if err := h.service.Delete(ctx, uid, id); err != nil {
		status := http.StatusInternalServerError
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
