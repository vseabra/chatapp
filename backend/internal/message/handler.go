package message

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
	group := r.Group(constants.APIv1 + "/rooms")
	group.Use(auth.AuthMiddleware(jwtSecret))
	group.GET(":id/messages", h.list)
}

func (h *Handler) list(c *gin.Context) {
	roomID := c.Param("id")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	cursor := c.Query("cursor")

	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.TIMEOUT_SECONDS)
	defer cancel()
	items, next, err := h.service.List(ctx, roomID, limit, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "nextCursor": next})
}
