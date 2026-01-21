package mysql

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql/model"
	"gorm.io/gorm"
)

// actionPassRepository 行動跳過儲存庫實作
type actionPassRepository struct {
	db *gorm.DB
}

// NewActionPassRepository 建立行動跳過儲存庫
func NewActionPassRepository(db *gorm.DB) repository.ActionPassRepository {
	return &actionPassRepository{db: db}
}

// CreateActionPass 建立跳過記錄
func (r *actionPassRepository) CreateActionPass(ctx context.Context, actionPass *entity.ActionPass) (*entity.ActionPass, error) {
	m := model.ActionPassModelFromEntity(actionPass)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetActionPassesByGameIDAndRound 取得某遊戲某回合的所有跳過記錄
func (r *actionPassRepository) GetActionPassesByGameIDAndRound(ctx context.Context, gameID int, round int) ([]*entity.ActionPass, error) {
	var models []model.ActionPassModel
	if err := r.db.WithContext(ctx).Where("game_id = ? AND round = ?", gameID, round).Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]*entity.ActionPass, 0, len(models))
	for _, m := range models {
		entities = append(entities, m.ToEntity())
	}
	return entities, nil
}

// CountActionPassesByGameIDAndRound 計算某遊戲某回合的跳過次數
func (r *actionPassRepository) CountActionPassesByGameIDAndRound(ctx context.Context, gameID int, round int) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.ActionPassModel{}).Where("game_id = ? AND round = ?", gameID, round).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// DeleteActionPassesByGameIDAndRound 刪除某遊戲某回合的所有跳過記錄
func (r *actionPassRepository) DeleteActionPassesByGameIDAndRound(ctx context.Context, gameID int, round int) error {
	return r.db.WithContext(ctx).Where("game_id = ? AND round = ?", gameID, round).Delete(&model.ActionPassModel{}).Error
}

// ExistsByGameIDPlayerIDAndRound 檢查某玩家是否已經跳過
func (r *actionPassRepository) ExistsByGameIDPlayerIDAndRound(ctx context.Context, gameID int, playerID int, round int) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.ActionPassModel{}).Where("game_id = ? AND player_id = ? AND round = ?", gameID, playerID, round).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
