package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/request"
	"github.com/Game-as-a-Service/The-Message/internal/adapter/sse"
	"github.com/Game-as-a-Service/The-Message/internal/usecase"
	"github.com/gin-gonic/gin"
)

// GameHandler 遊戲處理器
type GameHandler struct {
	gameUseCase   usecase.GameUseCase
	playerUseCase usecase.PlayerUseCase
	SSE           *sse.Event
}

// GameHandlerOptions 遊戲處理器選項
type GameHandlerOptions struct {
	Engine        *gin.Engine
	GameUseCase   usecase.GameUseCase
	PlayerUseCase usecase.PlayerUseCase
	SSE           *sse.Event
}

// RegisterGameHandler 註冊遊戲處理器
func RegisterGameHandler(opts *GameHandlerOptions) {
	handler := &GameHandler{
		gameUseCase:   opts.GameUseCase,
		playerUseCase: opts.PlayerUseCase,
		SSE:           opts.SSE,
	}

	opts.Engine.POST("/api/v1/games/start-legacy", handler.StartGame)
	opts.Engine.GET("/api/v1/games/:gameId/events", sse.HeadersMiddleware(), opts.SSE.ServeHTTP(), handler.GameEvent)
}

// StartGame godoc
// @Summary Start a new game
// @Description Start a new game with specified players
// @Tags games
// @Accept json
// @Produce json
// @Param players body request.CreateGameRequest true "Players"
// @Success 200 {object} response.CreateGameResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/games/start-legacy [post]
func (h *GameHandler) StartGame(c *gin.Context) {
	var req request.CreateGameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 轉換請求格式
	createReq := usecase.CreateGameRequest{
		Players: make([]usecase.PlayerInfo, len(req.Players)),
	}
	for i, p := range req.Players {
		createReq.Players[i] = usecase.PlayerInfo{
			ID:   p.ID,
			Name: p.Name,
		}
	}

	game, err := h.gameUseCase.StartGame(c, createReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Id":    game.ID,
		"Token": game.Token,
	})
}

// GetGame godoc
// @Summary Get game by ID
// @Description Get game details by game ID
// @Tags games
// @Produce json
// @Param gameId path int true "Game ID"
// @Success 200 {object} response.CreateGameResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId} [get]
func (h *GameHandler) GetGame(c *gin.Context) {
	gameID, _ := strconv.Atoi(c.Param("gameId"))

	game, err := h.gameUseCase.GetGameById(c, gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Id":    game.ID,
		"Token": game.Token,
	})
}

// DeleteGame godoc
// @Summary Delete a game
// @Description Delete a game by game ID
// @Tags games
// @Produce json
// @Param gameId path int true "Game ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId} [delete]
func (h *GameHandler) DeleteGame(c *gin.Context) {
	gameID, _ := strconv.Atoi(c.Param("gameId"))

	err := h.gameUseCase.DeleteGame(c, gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Game deleted"})
}

// GameSSERequest 遊戲 SSE 請求
type GameSSERequest struct {
	GameID  int    `json:"game_id,int"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// GameEvent godoc
// @Summary Get game events
// @Description Get game events via SSE
// @Tags games
// @Accept json
// @Produce json
// @Param gameId path int true "Game ID"
// @Success 200 {object} GameSSERequest
// @Router /api/v1/games/{gameId}/events [get]
func (h *GameHandler) GameEvent(c *gin.Context) {
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		return
	}

	v, ok := c.Get("clientChan")
	if !ok {
		return
	}

	clientChan, ok := v.(sse.ClientChan)
	if !ok {
		return
	}

	game, err := h.gameUseCase.GetGameById(c, gameID)
	if err != nil {
		return
	}

	h.SSE.Message <- gin.H{
		"message":        game.Status,
		"status":         game.Status,
		"game_id":        gameID,
		"current_player": game.CurrentPlayerID,
	}

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-clientChan; ok {
			log.Printf("msg: %+v", msg)
			data := GameSSERequest{}
			err := json.Unmarshal([]byte(msg), &data)
			if err != nil {
				log.Fatalf(err.Error())
			}

			if data.GameID == gameID {
				c.SSEvent("message", msg)
			}
			return true
		}
		return false
	})
}
