package repository

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// GameRepository 遊戲倉儲介面
// 定義於 domain 層，由 infrastructure 層實作
type GameRepository interface {
	// CreateGame 建立遊戲
	CreateGame(ctx context.Context, game *entity.Game) (*entity.Game, error)

	// GetGameById 根據 ID 查詢遊戲
	GetGameById(ctx context.Context, id int) (*entity.Game, error)

	// GetGameWithPlayers 查詢遊戲及其玩家
	GetGameWithPlayers(ctx context.Context, id int) (*entity.Game, error)

	// UpdateGame 更新遊戲
	UpdateGame(ctx context.Context, game *entity.Game) error

	// DeleteGame 刪除遊戲
	DeleteGame(ctx context.Context, id int) error
}
