package handlers

import (
	"net/http"

	"github.com/LiquidCats/upgrader/configs"
	"github.com/gin-gonic/gin"
)

type VersionHandler struct {
	cfg configs.Workers
}

func NewApiHandler(cfg configs.Workers) *VersionHandler {
	return &VersionHandler{cfg: cfg}
}

func (h *VersionHandler) Handle(c *gin.Context) {
	var endpoints []string

	for _, cfg := range h.cfg {
		endpoints = append(endpoints, cfg.ToWebsocket)
	}

	c.JSON(http.StatusOK, gin.H{"endpoints": endpoints})
}
