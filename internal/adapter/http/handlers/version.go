package handlers

import (
	"net/http"

	"github.com/LiquidCats/upgrader/configs"
	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	cfg configs.WorkersConfig
}

func NewAPIHandler(cfg configs.WorkersConfig) *APIHandler {
	return &APIHandler{cfg: cfg}
}

func (h *APIHandler) Handle(c *gin.Context) {
	var endpoints []string

	for _, cfg := range h.cfg {
		endpoints = append(endpoints, cfg.ToWebsocket)
	}

	c.JSON(http.StatusOK, gin.H{"endpoints": endpoints})
}
