package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RootHandler struct {
}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

func (h *RootHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "v1"})
}
