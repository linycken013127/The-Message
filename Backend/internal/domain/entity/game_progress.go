package entity

import "time"

// GameProgress 遊戲進度領域實體
type GameProgress struct {
	ID             int
	PlayerID       int
	GameID         int
	CardID         int
	Action         string
	TargetPlayerID int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewGameProgress 建立新遊戲進度（工廠方法）
func NewGameProgress(playerID int, gameID int, cardID int, action string, targetPlayerID int) *GameProgress {
	now := time.Now()
	return &GameProgress{
		PlayerID:       playerID,
		GameID:         gameID,
		CardID:         cardID,
		Action:         action,
		TargetPlayerID: targetPlayerID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// 遊戲動作常數
const (
	ActionTransmitIntelligence = "傳情報"
)
