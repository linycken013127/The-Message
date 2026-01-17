package entity

import "time"

// Game 遊戲領域實體
// 純 Go 結構體，不依賴任何外部框架
type Game struct {
	ID              int
	Token           string
	Status          string
	CurrentPlayerID int
	Players         []Player
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewGame 建立新遊戲（工廠方法）
func NewGame(token string, status string) *Game {
	now := time.Now()
	return &Game{
		Token:     token,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetCurrentPlayer 設定當前玩家
func (g *Game) SetCurrentPlayer(playerID int) {
	g.CurrentPlayerID = playerID
	g.UpdatedAt = time.Now()
}

// SetStatus 設定遊戲狀態
func (g *Game) SetStatus(status string) {
	g.Status = status
	g.UpdatedAt = time.Now()
}

// IsEnded 檢查遊戲是否結束
func (g *Game) IsEnded() bool {
	return g.Status == GameStatusEnd
}

// 遊戲狀態常數
const (
	GameStatusStart                 = "開始遊戲"
	GameStatusActionCardStage       = "功能牌階段"
	GameStatusTransmitIntelligence  = "情報牌階段"
	GameStatusEnd                   = "結束遊戲"
)
