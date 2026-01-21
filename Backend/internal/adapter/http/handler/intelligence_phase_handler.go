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
	opts.Engine.POST("/api/v1/games/:gameId/intelligence/accept", handler.AcceptIntelligence)
	opts.Engine.POST("/api/v1/games/:gameId/intelligence/reject", handler.RejectIntelligence)
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

// AcceptIntelligenceRequest 接收情報請求
type AcceptIntelligenceRequest struct {
	PlayerID int `json:"player_id" binding:"required"`
}

// AcceptIntelligence godoc
// @Summary Accept intelligence
// @Description Current target player accepts the intelligence card
// @Tags intelligence-phase
// @Accept json
// @Produce json
// @Param gameId path int true "Game ID"
// @Param request body AcceptIntelligenceRequest true "Accept intelligence request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId}/intelligence/accept [post]
func (h *IntelligencePhaseHandler) AcceptIntelligence(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "無效的遊戲 ID")
		return
	}

	var req AcceptIntelligenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "請求格式錯誤")
		return
	}

	result, err := h.intelligencePhaseUseCase.AcceptIntelligence(c, gameID, req.PlayerID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "情報已接收",
		"player_id": result.PlayerID,
		"card_id":   result.CardID,
	})
}

// RejectIntelligenceRequest 拒絕情報請求
type RejectIntelligenceRequest struct {
	PlayerID int `json:"player_id" binding:"required"`
}

// RejectIntelligence godoc
// @Summary Reject intelligence
// @Description Current target player rejects the intelligence card
// @Tags intelligence-phase
// @Accept json
// @Produce json
// @Param gameId path int true "Game ID"
// @Param request body RejectIntelligenceRequest true "Reject intelligence request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId}/intelligence/reject [post]
func (h *IntelligencePhaseHandler) RejectIntelligence(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "無效的遊戲 ID")
		return
	}

	var req RejectIntelligenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "請求格式錯誤")
		return
	}

	result, err := h.intelligencePhaseUseCase.RejectIntelligence(c, gameID, req.PlayerID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 根據是否自動接收回傳不同訊息
	if result.AutoAccepted {
		c.JSON(http.StatusOK, gin.H{
			"success":          true,
			"message":          result.Message,
			"auto_accepted":    true,
			"auto_accepted_by": result.AutoAcceptedBy,
		})
		return
	}

	// 情報傳給下一位玩家
	resp := gin.H{
		"success": true,
		"message": result.Message,
		"transfer": gin.H{
			"id":                       result.Transfer.ID,
			"game_id":                  result.Transfer.GameID,
			"card_id":                  result.Transfer.CardID,
			"sender_player_id":         result.Transfer.SenderPlayerID,
			"current_target_player_id": result.Transfer.CurrentTargetPlayerID,
			"face_up":                  result.Transfer.FaceUp,
			"status":                   result.Transfer.Status,
		},
	}

	// 如果明牌，則顯示卡片資訊
	if result.Transfer.FaceUp && result.Transfer.Card != nil {
		resp["card"] = gin.H{
			"id":    result.Transfer.Card.ID,
			"name":  result.Transfer.Card.Name,
			"color": result.Transfer.Card.Color,
		}
	}

	c.JSON(http.StatusOK, resp)
}
