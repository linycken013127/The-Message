package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HeartbeatHandler 心跳處理器
type HeartbeatHandler struct {
	Engine *gin.Engine
}

// RegisterHeartbeatHandler 註冊心跳處理器
func RegisterHeartbeatHandler(opts *HeartbeatHandler) {
	handler := &HeartbeatHandler{}

	opts.Engine.GET("/api/v1/heartbeat", handler.Heartbeat)
}

// Heartbeat godoc
// @Summary Check if the server is alive
// @Description Check if the server is alive
// @Tags heartbeat
// @Accept json
// @Produce json
// @Success 204
// @Router /api/v1/heartbeat [GET]
func (h *HeartbeatHandler) Heartbeat(c *gin.Context) {
	c.JSON(http.StatusNoContent, http.NoBody)
}
