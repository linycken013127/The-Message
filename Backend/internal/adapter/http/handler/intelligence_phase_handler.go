package handler

import (
	"net/http"
	"strconv"

	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/response"
	"github.com/Game-as-a-Service/The-Message/internal/usecase"
	"github.com/gin-gonic/gin"
)

// IntelligencePhaseHandler 情報階段處理器
type IntelligencePhaseHandler struct {
	intelligencePhaseUseCase usecase.IntelligencePhaseUseCase
}

// IntelligencePhaseHandlerOptions 情報階段處理器選項
type IntelligencePhaseHandlerOptions struct {
	Engine                   *gin.Engine
	IntelligencePhaseUseCase usecase.IntelligencePhaseUseCase
}

// RegisterIntelligencePhaseHandler 註冊情報階段處理器
func RegisterIntelligencePhaseHandler(opts *IntelligencePhaseHandlerOptions) {
	handler := &IntelligencePhaseHandler{
		intelligencePhaseUseCase: opts.IntelligencePhaseUseCase,
	}

	opts.Engine.POST("/api/v1/games/:gameId/intelligence/pass-card", handler.PassIntelligenceCard)
	opts.Engine.GET("/api/v1/games/:gameId/intelligence/active", handler.GetActiveTransfer)
}

// PassIntelligenceCardRequest 傳遞情報牌請求
type PassIntelligenceCardRequest struct {
	PlayerID       int `json:"player_id" binding:"required"`
	CardID         int `json:"card_id" binding:"required"`
	TargetPlayerID int `json:"target_player_id"` // 直達卡牌需要指定目標
}

// PassIntelligenceCard godoc
// @Summary Pass intelligence card
// @Description Current player passes an intelligence card
// @Tags intelligence-phase
// @Accept json
// @Produce json
// @Param gameId path int true "Game ID"
// @Param request body PassIntelligenceCardRequest true "Pass intelligence card request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId}/intelligence/pass-card [post]
func (h *IntelligencePhaseHandler) PassIntelligenceCard(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "無效的遊戲 ID")
		return
	}

	var req PassIntelligenceCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "請求格式錯誤")
		return
	}

	transfer, err := h.intelligencePhaseUseCase.PassIntelligenceCard(c, gameID, req.PlayerID, req.CardID, req.TargetPlayerID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 建立回應
	resp := gin.H{
		"success": true,
		"message": "情報牌已送出",
		"transfer": gin.H{
			"id":                       transfer.ID,
			"game_id":                  transfer.GameID,
			"card_id":                  transfer.CardID,
			"sender_player_id":         transfer.SenderPlayerID,
			"current_target_player_id": transfer.CurrentTargetPlayerID,
			"face_up":                  transfer.FaceUp,
			"status":                   transfer.Status,
		},
	}

	// 如果明牌，則顯示卡片資訊
	if transfer.FaceUp && transfer.Card != nil {
		resp["card"] = gin.H{
			"id":    transfer.Card.ID,
			"name":  transfer.Card.Name,
			"color": transfer.Card.Color,
		}
	}

	c.JSON(http.StatusOK, resp)
}

// GetActiveTransfer godoc
// @Summary Get active intelligence transfer
// @Description Get the currently active intelligence transfer in the game
// @Tags intelligence-phase
// @Accept json
// @Produce json
// @Param gameId path int true "Game ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId}/intelligence/active [get]
func (h *IntelligencePhaseHandler) GetActiveTransfer(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "無效的遊戲 ID")
		return
	}

	transfer, err := h.intelligencePhaseUseCase.GetActiveTransfer(c, gameID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if transfer == nil {
		c.JSON(http.StatusOK, gin.H{
			"has_active_transfer": false,
		})
		return
	}

	resp := gin.H{
		"has_active_transfer": true,
		"transfer": gin.H{
			"id":                       transfer.ID,
			"game_id":                  transfer.GameID,
			"card_id":                  transfer.CardID,
			"sender_player_id":         transfer.SenderPlayerID,
			"current_target_player_id": transfer.CurrentTargetPlayerID,
			"face_up":                  transfer.FaceUp,
			"status":                   transfer.Status,
		},
	}

	// 如果明牌，則顯示卡片資訊
	if transfer.FaceUp && transfer.Card != nil {
		resp["card"] = gin.H{
			"id":    transfer.Card.ID,
			"name":  transfer.Card.Name,
			"color": transfer.Card.Color,
		}
	}

	c.JSON(http.StatusOK, resp)
}
