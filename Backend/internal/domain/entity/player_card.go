package entity

import "time"

// PlayerCard 玩家手牌領域實體
type PlayerCard struct {
	ID        int
	PlayerID  int
	GameID    int
	CardID    int
	Type      string
	Card      Card
	Player    Player
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewPlayerCard 建立新玩家手牌（工廠方法）
func NewPlayerCard(playerID int, gameID int, cardID int, cardType string) *PlayerCard {
	now := time.Now()
	return &PlayerCard{
		PlayerID:  playerID,
		GameID:    gameID,
		CardID:    cardID,
		Type:      cardType,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// 手牌類型常數
const (
	PlayerCardTypeHand         = "hand"         // 手牌
	PlayerCardTypeIntelligence = "intelligence" // 情報
)
