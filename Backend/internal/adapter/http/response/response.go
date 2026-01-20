package response

import "github.com/gin-gonic/gin"

// CreateGameResponse 建立遊戲回應
type CreateGameResponse struct {
	ID    int    `json:"Id"`
	Token string `json:"Token"`
}

// PlayCardResponse 出牌回應
type PlayCardResponse struct {
	Result bool `json:"result"`
}

// PlayerCardsResponse 玩家手牌回應
type PlayerCardsResponse struct {
	PlayerCards []map[string]interface{} `json:"player_cards"`
}

// ErrorResponse 錯誤回應
type ErrorResponse struct {
	Message string `json:"message"`
}

// RegisterPlayerResponse 註冊玩家回應
type RegisterPlayerResponse struct {
	PlayerID   int    `json:"playerId"`
	PlayerName string `json:"playerName"`
}

// Success 成功回應
func Success(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

// Error 錯誤回應
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}
