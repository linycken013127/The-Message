package handler

import (
	"net/http"
	"strconv"

	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/middleware"
	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/response"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/usecase"
	"github.com/gin-gonic/gin"
)

// GameQueryHandler 遊戲查詢處理器
type GameQueryHandler struct {
	gameQueryUseCase usecase.GameQueryUseCase
}

// GameQueryHandlerOptions 遊戲查詢處理器選項
type GameQueryHandlerOptions struct {
	Engine           *gin.Engine
	GameQueryUseCase usecase.GameQueryUseCase
	AccountRepo      repository.AccountRepository
}

// RegisterGameQueryHandler 註冊遊戲查詢處理器
func RegisterGameQueryHandler(opts *GameQueryHandlerOptions) {
	handler := &GameQueryHandler{
		gameQueryUseCase: opts.GameQueryUseCase,
	}

	// 需要認證的路由
	authGroup := opts.Engine.Group("/api/v1")
	authGroup.Use(middleware.AuthMiddleware(opts.AccountRepo))
	{
		authGroup.GET("/games/:gameId", handler.GetGameInfo)
	}
}

// GetGameInfo godoc
// @Summary Get game information
// @Description Get game information with personal and public player info
// @Tags game
// @Accept json
// @Produce json
// @Param X-Account-ID header int true "Account ID"
// @Param gameId path int true "Game ID"
// @Success 200 {object} usecase.GameInfoResult
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/games/{gameId} [get]
func (h *GameQueryHandler) GetGameInfo(c *gin.Context) {
	accountID := middleware.GetAccountID(c)
	gameID, err := strconv.Atoi(c.Param("gameId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "無效的遊戲 ID")
		return
	}

	result, err := h.gameQueryUseCase.GetGameInfo(c, gameID, accountID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}

	response.Success(c, http.StatusOK, result)
}
