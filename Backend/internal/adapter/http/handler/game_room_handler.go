package handler

import (
	"net/http"
	"strconv"

	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/middleware"
	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/request"
	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/response"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/usecase"
	"github.com/gin-gonic/gin"
)

// GameRoomHandler 遊戲房處理器
type GameRoomHandler struct {
	gameRoomUseCase usecase.GameRoomUseCase
}

// GameRoomHandlerOptions 遊戲房處理器選項
type GameRoomHandlerOptions struct {
	Engine          *gin.Engine
	GameRoomUseCase usecase.GameRoomUseCase
	AccountRepo     repository.AccountRepository
}

// RegisterGameRoomHandler 註冊遊戲房處理器
func RegisterGameRoomHandler(opts *GameRoomHandlerOptions) {
	handler := &GameRoomHandler{
		gameRoomUseCase: opts.GameRoomUseCase,
	}

	// 需要認證的路由
	authGroup := opts.Engine.Group("/api/v1")
	authGroup.Use(middleware.AuthMiddleware(opts.AccountRepo))
	{
		authGroup.POST("/games", handler.CreateGameRoom)
		authGroup.POST("/games/:gameId/join", handler.JoinGameRoom)
		authGroup.POST("/games/:gameId/start", handler.StartGame)
	}
}

// CreateGameRoom godoc
// @Summary Create a new game room
// @Description Create a new game room with the authenticated user as host
// @Tags game-room
// @Accept json
// @Produce json
// @Param X-Account-ID header int true "Account ID"
// @Param request body request.CreateGameRoomRequest true "Create game room request"
// @Success 200 {object} response.CreateGameRoomResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/v1/games [post]
func (h *GameRoomHandler) CreateGameRoom(c *gin.Context) {
	accountID := middleware.GetAccountID(c)

	var req request.CreateGameRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果沒有 body 也可以，使用預設值
		req.MaxPlayers = 9
	}

	game, err := h.gameRoomUseCase.CreateGameRoom(c, accountID, req.MaxPlayers)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, response.CreateGameRoomResponse{
		GameID:       game.ID,
		Status:       game.Status,
		HostPlayerID: game.HostAccountID,
	})
}

// JoinGameRoom godoc
// @Summary Join a game room
// @Description Join an existing game room
// @Tags game-room
// @Accept json
// @Produce json
// @Param X-Account-ID header int true "Account ID"
// @Param gameId path int true "Game ID"
// @Success 200 {object} response.JoinGameRoomResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId}/join [post]
func (h *GameRoomHandler) JoinGameRoom(c *gin.Context) {
	accountID := middleware.GetAccountID(c)
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "無效的遊戲 ID")
		return
	}

	err = h.gameRoomUseCase.JoinGameRoom(c, gameID, accountID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, response.JoinGameRoomResponse{
		GameID:   gameID,
		PlayerID: accountID,
		Message:  "加入成功",
	})
}

// StartGame godoc
// @Summary Start a game
// @Description Start a game (only host can start)
// @Tags game-room
// @Accept json
// @Produce json
// @Param X-Account-ID header int true "Account ID"
// @Param gameId path int true "Game ID"
// @Success 200 {object} response.StartGameResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId}/start [post]
func (h *GameRoomHandler) StartGame(c *gin.Context) {
	accountID := middleware.GetAccountID(c)
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "無效的遊戲 ID")
		return
	}

	game, err := h.gameRoomUseCase.StartGame(c, gameID, accountID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, response.StartGameResponse{
		GameID: game.ID,
		Status: game.Status,
	})
}
