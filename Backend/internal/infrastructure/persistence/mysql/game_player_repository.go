package mysql

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql/model"
	"gorm.io/gorm"
)

// GamePlayerRepository MySQL 遊戲玩家關聯倉儲實作
type GamePlayerRepository struct {
	db *gorm.DB
}

// NewGamePlayerRepository 建立遊戲玩家關聯倉儲
func NewGamePlayerRepository(db *gorm.DB) repository.GamePlayerRepository {
	return &GamePlayerRepository{
		db: db,
	}
}

// CreateGamePlayer 建立遊戲玩家關聯
func (r *GamePlayerRepository) CreateGamePlayer(ctx context.Context, gamePlayer *entity.GamePlayer) (*entity.GamePlayer, error) {
	gamePlayerModel := model.GamePlayerModelFromEntity(gamePlayer)

	err := r.db.Create(gamePlayerModel).Error
	if err != nil {
		return nil, err
	}

	return gamePlayerModel.ToEntity(), nil
}

// GetGamePlayersByGameID 根據遊戲 ID 取得所有玩家關聯
func (r *GamePlayerRepository) GetGamePlayersByGameID(ctx context.Context, gameID int) ([]*entity.GamePlayer, error) {
	var gamePlayerModels []model.GamePlayerModel

	result := r.db.Where("game_id = ?", gameID).Order("join_order ASC").Find(&gamePlayerModels)
	if result.Error != nil {
		return nil, result.Error
	}

	gamePlayers := make([]*entity.GamePlayer, len(gamePlayerModels))
	for i, gpm := range gamePlayerModels {
		gamePlayers[i] = gpm.ToEntity()
	}

	return gamePlayers, nil
}

// ExistsByGameIDAndAccountID 檢查玩家是否已在遊戲中
func (r *GamePlayerRepository) ExistsByGameIDAndAccountID(ctx context.Context, gameID int, accountID int) (bool, error) {
	var count int64

	result := r.db.Model(&model.GamePlayerModel{}).
		Where("game_id = ? AND account_id = ?", gameID, accountID).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

// CountByGameID 計算遊戲中的玩家數量
func (r *GamePlayerRepository) CountByGameID(ctx context.Context, gameID int) (int, error) {
	var count int64

	result := r.db.Model(&model.GamePlayerModel{}).
		Where("game_id = ?", gameID).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

// GetGamePlayerByGameIDAndAccountID 根據遊戲 ID 和帳號 ID 取得玩家關聯
func (r *GamePlayerRepository) GetGamePlayerByGameIDAndAccountID(ctx context.Context, gameID int, accountID int) (*entity.GamePlayer, error) {
	var gamePlayerModel model.GamePlayerModel

	result := r.db.Where("game_id = ? AND account_id = ?", gameID, accountID).First(&gamePlayerModel)
	if result.Error != nil {
		return nil, result.Error
	}

	return gamePlayerModel.ToEntity(), nil
}
