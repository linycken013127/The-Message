package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/request"
	"github.com/Game-as-a-Service/The-Message/internal/adapter/sse"
	"github.com/Game-as-a-Service/The-Message/internal/usecase"
	"github.com/gin-gonic/gin"
)

// PlayerHandler 玩家處理器
type PlayerHandler struct {
	playerUseCase usecase.PlayerUseCase
	gameUseCase   usecase.GameUseCase
	SSE           *sse.Event
}

// PlayerHandlerOptions 玩家處理器選項
type PlayerHandlerOptions struct {
	Engine        *gin.Engine
	PlayerUseCase usecase.PlayerUseCase
	GameUseCase   usecase.GameUseCase
	SSE           *sse.Event
}

// RegisterPlayerHandler 註冊玩家處理器
func RegisterPlayerHandler(opts *PlayerHandlerOptions) {
	handler := &PlayerHandler{
		playerUseCase: opts.PlayerUseCase,
		gameUseCase:   opts.GameUseCase,
		SSE:           opts.SSE,
	}

	opts.Engine.POST("/api/v1/players/:playerId/player-cards", handler.PlayCard)
	opts.Engine.POST("/api/v1/player/:playerId/transmit-intelligence", handler.TransmitIntelligence)
	opts.Engine.POST("/api/v1/players/:playerId/accept", handler.AcceptCard)
}

// PlayCard godoc
// @Summary Play a card
// @Description Play a card from player's hand
// @Tags players
// @Accept json
// @Produce json
// @Param playerId path int true "Player ID"
// @Param card_id body request.PlayCardRequest true "Card ID"
// @Success 200 {object} response.PlayCardResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/players/{playerId}/player-cards [post]
func (h *PlayerHandler) PlayCard(c *gin.Context) {
	playerID, _ := strconv.Atoi(c.Param("playerId"))
	var req request.PlayCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, card, err := h.playerUseCase.PlayCard(c, playerID, req.CardID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	h.SSE.Message <- gin.H{
		"game_id":     game.ID,
		"status":      game.Status,
		"message":     fmt.Sprintf("玩家: %d 已出牌", playerID),
		"card":        card.Name,
		"next_player": game.CurrentPlayerID,
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}

// TransmitIntelligence godoc
// @Summary Transmit intelligence
// @Description Transmit an intelligence card
// @Tags players
// @Accept json
// @Produce json
// @Param playerId path int true "Player ID"
// @Param card_id body request.PlayCardRequest true "Card ID"
// @Success 200 {object} response.PlayCardResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/player/{playerId}/transmit-intelligence [post]
func (h *PlayerHandler) TransmitIntelligence(c *gin.Context) {
	playerID, _ := strconv.Atoi(c.Param("playerId"))
	var req request.PlayCardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	player, err := h.playerUseCase.GetPlayerById(c, playerID)
	if err != nil || player == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Player not found"})
		return
	}

	exist, err := h.playerUseCase.CheckPlayerCardExist(c, playerID, player.GameID, req.CardID)
	if err != nil || !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Card not found"})
		return
	}

	ret, err := h.playerUseCase.TransmitIntelligenceCard(c, playerID, player.GameID, req.CardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": ret,
	})
}

// AcceptCard godoc
// @Summary Accept Card
// @Description Decide accept card or not
// @Tags players
// @Accept json
// @Produce json
// @Param playerId path int true "Player ID"
// @Param accept body request.AcceptCardRequest true "Accept"
// @Success 200 {object} response.PlayCardResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/players/{playerId}/accept [post]
func (h *PlayerHandler) AcceptCard(c *gin.Context) {
	playerID, _ := strconv.Atoi(c.Param("playerId"))
	var req request.AcceptCardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := h.playerUseCase.AcceptCard(c, playerID, req.Accept)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	winner, err := h.playerUseCase.CheckWin(c, playerID)
	if winner != nil {
		h.SSE.Message <- gin.H{
			"game_id": winner.Game.ID,
			"status":  winner.Game.Status,
			"message": fmt.Sprintf("玩家: %d 已贏得遊戲", playerID),
			"winner":  winner.Name,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}
