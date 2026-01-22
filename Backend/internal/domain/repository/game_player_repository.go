package repository

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// GamePlayerRepository 遊戲玩家關聯倉儲介面
type GamePlayerRepository interface {
	// CreateGamePlayer 建立遊戲玩家關聯
	CreateGamePlayer(ctx context.Context, gamePlayer *entity.GamePlayer) (*entity.GamePlayer, error)

	// GetGamePlayersByGameID 根據遊戲 ID 取得所有玩家關聯
	GetGamePlayersByGameID(ctx context.Context, gameID int) ([]*entity.GamePlayer, error)

	// ExistsByGameIDAndAccountID 檢查玩家是否已在遊戲中
	ExistsByGameIDAndAccountID(ctx context.Context, gameID int, accountID int) (bool, error)

	// CountByGameID 計算遊戲中的玩家數量
	CountByGameID(ctx context.Context, gameID int) (int, error)

	// GetGamePlayerByGameIDAndAccountID 根據遊戲 ID 和帳號 ID 取得玩家關聯
	GetGamePlayerByGameIDAndAccountID(ctx context.Context, gameID int, accountID int) (*entity.GamePlayer, error)
}
