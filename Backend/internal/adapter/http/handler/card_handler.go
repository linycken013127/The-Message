package handler

import (
	"net/http"
	"strconv"

	"github.com/Game-as-a-Service/The-Message/internal/usecase"
	"github.com/gin-gonic/gin"
)

// CardHandler 卡片處理器
type CardHandler struct {
	cardUseCase usecase.CardUseCase
}

// CardHandlerOptions 卡片處理器選項
type CardHandlerOptions struct {
	Engine      *gin.Engine
	CardUseCase usecase.CardUseCase
}

// RegisterCardHandler 註冊卡片處理器
func RegisterCardHandler(opts *CardHandlerOptions) {
	handler := &CardHandler{
		cardUseCase: opts.CardUseCase,
	}

	opts.Engine.GET("/api/v1/player/:playerId/player-cards/", handler.GetPlayerCards)
}

// GetPlayerCards godoc
// @Summary GetPlayerCards
// @Description Get player's cards by player ID
// @Tags player_cards
// @Produce json
// @Param playerId path int true "Player ID"
// @Success 200 {object} response.PlayerCardsResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/player/{playerId}/player-cards/ [get]
func (h *CardHandler) GetPlayerCards(c *gin.Context) {
	playerID, _ := strconv.Atoi(c.Param("playerId"))
	playerCards, err := h.cardUseCase.GetPlayerCardsByPlayerId(c, playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	playerCardsInfo := []map[string]interface{}{}
	for _, card := range playerCards {
		dict := map[string]interface{}{
			"id":    card.ID,
			"name":  card.Name,
			"color": card.Color,
		}
		playerCardsInfo = append(playerCardsInfo, dict)
	}

	c.JSON(http.StatusOK, gin.H{"player_cards": playerCardsInfo})
}
