package entity

import "time"

// ActionPass 行動階段跳過記錄
type ActionPass struct {
	ID        int
	GameID    int
	PlayerID  int
	Round     int
	CreatedAt time.Time
}

// NewActionPass 建立新的跳過記錄
func NewActionPass(gameID int, playerID int, round int) *ActionPass {
	return &ActionPass{
		GameID:    gameID,
		PlayerID:  playerID,
		Round:     round,
		CreatedAt: time.Now(),
	}
}
