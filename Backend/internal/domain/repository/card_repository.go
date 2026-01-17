package repository

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// CardRepository 卡片倉儲介面
type CardRepository interface {
	// GetCardById 根據 ID 查詢卡片
	GetCardById(ctx context.Context, id int) (*entity.Card, error)

	// CreateCard 建立卡片
	CreateCard(ctx context.Context, card *entity.Card) (*entity.Card, error)

	// GetCards 取得所有卡片
	GetCards(ctx context.Context) ([]*entity.Card, error)

	// GetPlayerCardsByPlayerId 根據玩家 ID 取得玩家手牌
	GetPlayerCardsByPlayerId(ctx context.Context, playerID int, gameID int) ([]*entity.Card, error)
}
