package entity

import "time"

// GamePlayer 遊戲玩家關聯實體
// 追蹤哪些帳號加入了哪些遊戲房
type GamePlayer struct {
	ID        int
	GameID    int
	AccountID int
	JoinOrder int // 加入順序（用於決定玩家順序）
	CreatedAt time.Time
}

// NewGamePlayer 建立新的遊戲玩家關聯
func NewGamePlayer(gameID int, accountID int, joinOrder int) *GamePlayer {
	return &GamePlayer{
		GameID:    gameID,
		AccountID: accountID,
		JoinOrder: joinOrder,
		CreatedAt: time.Now(),
	}
}
