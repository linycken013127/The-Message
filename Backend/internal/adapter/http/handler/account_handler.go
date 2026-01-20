package handler

import (
	"net/http"

	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/request"
	"github.com/Game-as-a-Service/The-Message/internal/adapter/http/response"
	"github.com/Game-as-a-Service/The-Message/internal/usecase"
	"github.com/gin-gonic/gin"
)

// AccountHandler 帳號處理器
type AccountHandler struct {
	accountUseCase usecase.AccountUseCase
}

// AccountHandlerOptions 帳號處理器選項
type AccountHandlerOptions struct {
	Engine         *gin.Engine
	AccountUseCase usecase.AccountUseCase
}

// RegisterAccountHandler 註冊帳號處理器
func RegisterAccountHandler(opts *AccountHandlerOptions) {
	handler := &AccountHandler{
		accountUseCase: opts.AccountUseCase,
	}

	opts.Engine.POST("/api/v1/players", handler.RegisterPlayer)
}

// RegisterPlayer godoc
// @Summary Register a new player
// @Description Register a new player with a unique name
// @Tags players
// @Accept json
// @Produce json
// @Param request body request.RegisterPlayerRequest true "Player registration request"
// @Success 200 {object} response.RegisterPlayerResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/players [post]
func (h *AccountHandler) RegisterPlayer(c *gin.Context) {
	var req request.RegisterPlayerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "玩家名稱不可為空")
		return
	}

	account, err := h.accountUseCase.RegisterAccount(c, req.PlayerName)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusOK, response.RegisterPlayerResponse{
		PlayerID:   account.ID,
		PlayerName: account.Name,
	})
}
