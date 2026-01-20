package request

// CreateGameRequest 建立遊戲請求
type CreateGameRequest struct {
	Players []PlayerInfo `json:"players"`
}

// PlayerInfo 玩家資訊
type PlayerInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PlayCardRequest 出牌請求
type PlayCardRequest struct {
	CardID int `json:"card_id"`
}

// AcceptCardRequest 接收卡片請求
type AcceptCardRequest struct {
	Accept bool `json:"accept"`
}

// RegisterPlayerRequest 註冊玩家請求
type RegisterPlayerRequest struct {
	PlayerName string `json:"playerName" binding:"required"`
}

// CreateGameRoomRequest 建立遊戲房請求
type CreateGameRoomRequest struct {
	MaxPlayers int `json:"maxPlayers"`
}

// JoinGameRoomRequest 加入遊戲房請求
type JoinGameRoomRequest struct{}

// StartGameRequest 開始遊戲請求
type StartGameRequest struct{}
