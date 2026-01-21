package repository

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// ActionPassRepository 行動跳過儲存庫介面
type ActionPassRepository interface {
	// CreateActionPass 建立跳過記錄
	CreateActionPass(ctx context.Context, actionPass *entity.ActionPass) (*entity.ActionPass, error)
	// GetActionPassesByGameIDAndRound 取得某遊戲某回合的所有跳過記錄
	GetActionPassesByGameIDAndRound(ctx context.Context, gameID int, round int) ([]*entity.ActionPass, error)
	// CountActionPassesByGameIDAndRound 計算某遊戲某回合的跳過次數
	CountActionPassesByGameIDAndRound(ctx context.Context, gameID int, round int) (int, error)
	// DeleteActionPassesByGameIDAndRound 刪除某遊戲某回合的所有跳過記錄
	DeleteActionPassesByGameIDAndRound(ctx context.Context, gameID int, round int) error
	// ExistsByGameIDPlayerIDAndRound 檢查某玩家是否已經跳過
	ExistsByGameIDPlayerIDAndRound(ctx context.Context, gameID int, playerID int, round int) (bool, error)
}
