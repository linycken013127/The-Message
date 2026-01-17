package entity

import "time"

// Deck 牌組領域實體
type Deck struct {
	ID        int
	GameID    int
	CardID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewDeck 建立新牌組（工廠方法）
func NewDeck(gameID int, cardID int) *Deck {
	now := time.Now()
	return &Deck{
		GameID:    gameID,
		CardID:    cardID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
