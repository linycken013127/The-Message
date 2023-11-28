package http

import (
	"net/http"
	"strconv"

	"github.com/Game-as-a-Service/The-Message/service/service"
	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	playerService service.PlayerService
}

type PlayerHandlerOptions struct {
	Engine  *gin.Engine
	Service service.PlayerService
}

func RegisterPlayerHandler(opts *PlayerHandlerOptions) {
	handler := &PlayerHandler{
		playerService: opts.Service,
	}

	opts.Engine.GET("/api/v1/players/:playerId/player-cards", handler.GetPlayerCards)
	opts.Engine.POST("/api/v1/players/:playerId/player-cards/:cardId", handler.PlayCard)
}

func (p *PlayerHandler) PlayCard(c *gin.Context) {
	playerId, _ := strconv.Atoi(c.Param("playerId"))
	cardId, _ := strconv.Atoi(c.Param("cardId"))

	result, err := p.playerService.PlayCard(c, playerId, cardId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

func (p *PlayerHandler) GetPlayerCards(context *gin.Context) {
	playerId, _ := strconv.Atoi(context.Param("playerId"))

	playerCards, err := p.playerService.GetPlayerCards(context, playerId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var response []PlayerCardsResponse
	for _, playerCard := range *playerCards {
		response = append(response, PlayerCardsResponse{
			Id:    playerCard.Id,
			Name:  playerCard.Card.Name,
			Color: playerCard.Card.Color,
		})
	}

	context.JSON(http.StatusOK, gin.H{
		"hand_cards": response,
	})
}

type PlayerCardsResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
