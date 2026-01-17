package repository

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// GameProgressRepository 遊戲進度倉儲介面
type GameProgressRepository interface {
	// CreateGameProgress 建立遊戲進度
	CreateGameProgress(ctx context.Context, gameProgress *entity.GameProgress) (*entity.GameProgress, error)

	// GetGameProgress 取得遊戲進度
	GetGameProgress(ctx context.Context, targetPlayerID int, gameID int) (*entity.GameProgress, error)

	// UpdateGameProgress 更新遊戲進度
	UpdateGameProgress(ctx context.Context, gameProgress *entity.GameProgress, nextPlayerID int) (*entity.GameProgress, error)
}
