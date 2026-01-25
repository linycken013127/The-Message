package handler

import (
	"net/http"
	"strconv"

	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/response"
	"github.com/Game-as-a-Service/The-Message/internal/usecase"
	"github.com/gin-gonic/gin"
)

// ActionPhaseHandler 行動階段處理器
type ActionPhaseHandler struct {
	actionPhaseUseCase usecase.ActionPhaseUseCase
}

// ActionPhaseHandlerOptions 行動階段處理器選項
type ActionPhaseHandlerOptions struct {
	Engine             *gin.Engine
	ActionPhaseUseCase usecase.ActionPhaseUseCase
}

// RegisterActionPhaseHandler 註冊行動階段處理器
func RegisterActionPhaseHandler(opts *ActionPhaseHandlerOptions) {
	handler := &ActionPhaseHandler{
		actionPhaseUseCase: opts.ActionPhaseUseCase,
	}

	opts.Engine.POST("/api/v1/games/:gameId/actions/draw", handler.DrawCards)
	opts.Engine.POST("/api/v1/games/:gameId/actions/play-card", handler.PlayFunctionCard)
	opts.Engine.POST("/api/v1/games/:gameId/actions/pass", handler.Pass)
}

// DrawCardsRequest 抽牌請求
type DrawCardsRequest struct {
	PlayerID int `json:"player_id" binding:"required"`
}

// PlayCardRequest 出牌請求
type PlayCardRequest struct {
	PlayerID int `json:"player_id" binding:"required"`
	CardID   int `json:"card_id" binding:"required"`
}

// PassRequest 跳過請求
type PassRequest struct {
	PlayerID int `json:"player_id" binding:"required"`
}

// DrawCards godoc
// @Summary Draw cards
// @Description Current player draws 2 cards from the deck
// @Tags action-phase
// @Accept json
// @Produce json
// @Param gameId path int true "Game ID"
// @Param request body DrawCardsRequest true "Draw cards request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId}/actions/draw [post]
func (h *ActionPhaseHandler) DrawCards(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "無效的遊戲 ID")
		return
	}

	var req DrawCardsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "請求格式錯誤")
		return
	}

	cards, err := h.actionPhaseUseCase.DrawCards(c, gameID, req.PlayerID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	cardIDs := make([]int, len(cards))
	for i, card := range cards {
		cardIDs[i] = card.ID
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"card_ids": cardIDs,
		"count":    len(cards),
	})
}

// PlayFunctionCard godoc
// @Summary Play a function card
// @Description Current player plays a function card
// @Tags action-phase
// @Accept json
// @Produce json
// @Param gameId path int true "Game ID"
// @Param request body PlayCardRequest true "Play card request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId}/actions/play-card [post]
func (h *ActionPhaseHandler) PlayFunctionCard(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "無效的遊戲 ID")
		return
	}

	var req PlayCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "請求格式錯誤")
		return
	}

	card, err := h.actionPhaseUseCase.PlayFunctionCard(c, gameID, req.PlayerID, req.CardID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"card_id": card.ID,
	})
}

// Pass godoc
// @Summary Pass action
// @Description Current player passes their action
// @Tags action-phase
// @Accept json
// @Produce json
// @Param gameId path int true "Game ID"
// @Param request body PassRequest true "Pass request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId}/actions/pass [post]
func (h *ActionPhaseHandler) Pass(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "無效的遊戲 ID")
		return
	}

	var req PassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "請求格式錯誤")
		return
	}

	result, err := h.actionPhaseUseCase.Pass(c, gameID, req.PlayerID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":       true,
		"all_passed":    result.AllPassed,
		"phase_changed": result.AllPassed,
		"new_phase":     result.NewPhase,
	})
}
