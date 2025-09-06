package http

import (
	"github.com/gin-gonic/gin"
)

// NewRouter constructs a new Gin engine with default middleware.
func NewRouter() *gin.Engine {
	r := gin.Default()
	return r
}
