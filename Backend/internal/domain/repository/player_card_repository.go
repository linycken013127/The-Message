package repository

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// PlayerCardRepository 玩家手牌倉儲介面
type PlayerCardRepository interface {
	// GetPlayerCardById 根據 ID 查詢玩家手牌
	GetPlayerCardById(ctx context.Context, id int) (*entity.PlayerCard, error)

	// GetPlayerCardsByGameId 根據遊戲 ID 查詢玩家手牌
	GetPlayerCardsByGameId(ctx context.Context, gameID int) ([]*entity.PlayerCard, error)

	// CreatePlayerCard 建立玩家手牌
	CreatePlayerCard(ctx context.Context, card *entity.PlayerCard) (*entity.PlayerCard, error)

	// DeletePlayerCard 刪除玩家手牌
	DeletePlayerCard(ctx context.Context, id int) error

	// DeletePlayerCardByPlayerIdAndCardId 根據玩家 ID 和卡片 ID 刪除手牌
	DeletePlayerCardByPlayerIdAndCardId(ctx context.Context, playerID int, gameID int, cardID int) (bool, error)

	// ExistPlayerCardByPlayerIdAndCardId 檢查玩家手牌是否存在
	ExistPlayerCardByPlayerIdAndCardId(ctx context.Context, playerID int, gameID int, cardID int) (bool, error)

	// GetPlayerCards 取得玩家手牌
	GetPlayerCards(ctx context.Context, playerCard *entity.PlayerCard) (*[]entity.PlayerCard, error)

	// CountHandCardsByPlayerID 計算玩家手牌數量
	CountHandCardsByPlayerID(ctx context.Context, playerID int) (int, error)
}
