package repository

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// DeckRepository 牌組倉儲介面
type DeckRepository interface {
	// CreateDeck 建立牌組
	CreateDeck(ctx context.Context, deck *entity.Deck) (*entity.Deck, error)

	// GetDecksByGameId 根據遊戲 ID 查詢牌組
	GetDecksByGameId(ctx context.Context, gameID int) ([]*entity.Deck, error)

	// DeleteDeck 刪除牌組
	DeleteDeck(ctx context.Context, id int) error
}
