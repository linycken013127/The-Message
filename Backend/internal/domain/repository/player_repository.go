package repository

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// PlayerRepository 玩家倉儲介面
type PlayerRepository interface {
	// CreatePlayer 建立玩家
	CreatePlayer(ctx context.Context, player *entity.Player) (*entity.Player, error)

	// GetPlayerById 根據 ID 查詢玩家
	GetPlayerById(ctx context.Context, id int) (*entity.Player, error)

	// GetPlayersByGameId 根據遊戲 ID 查詢所有玩家
	GetPlayersByGameId(ctx context.Context, gameID int) ([]*entity.Player, error)

	// GetPlayerWithPlayerCards 查詢玩家及其手牌
	GetPlayerWithPlayerCards(ctx context.Context, playerID int) (*entity.Player, error)

	// GetPlayerWithGame 查詢玩家及其遊戲
	GetPlayerWithGame(ctx context.Context, playerID int) (*entity.Player, error)

	// GetPlayerWithGamePlayersAndPlayerCardsCard 查詢玩家完整關聯資料
	GetPlayerWithGamePlayersAndPlayerCardsCard(ctx context.Context, playerID int) (*entity.Player, error)
}
