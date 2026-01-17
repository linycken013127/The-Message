package mysql

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql/model"
	"gorm.io/gorm"
)

// GameProgressRepository MySQL 遊戲進度倉儲實作
type GameProgressRepository struct {
	db *gorm.DB
}

// NewGameProgressRepository 建立遊戲進度倉儲
func NewGameProgressRepository(db *gorm.DB) repository.GameProgressRepository {
	return &GameProgressRepository{
		db: db,
	}
}

// CreateGameProgress 建立遊戲進度
func (r *GameProgressRepository) CreateGameProgress(ctx context.Context, gameProgress *entity.GameProgress) (*entity.GameProgress, error) {
	progressModel := model.GameProgressModelFromEntity(gameProgress)

	result := r.db.Create(progressModel)
	if result.Error != nil {
		return nil, result.Error
	}

	return progressModel.ToEntity(), nil
}

// GetGameProgress 取得遊戲進度
func (r *GameProgressRepository) GetGameProgress(ctx context.Context, targetPlayerID int, gameID int) (*entity.GameProgress, error) {
	var progressModel model.GameProgressModel

	result := r.db.First(&progressModel, "target_player_id = ? AND game_id = ?", targetPlayerID, gameID)
	if result.Error != nil {
		return nil, result.Error
	}

	return progressModel.ToEntity(), nil
}

// UpdateGameProgress 更新遊戲進度
func (r *GameProgressRepository) UpdateGameProgress(ctx context.Context, gameProgress *entity.GameProgress, nextPlayerID int) (*entity.GameProgress, error) {
	progressModel := model.GameProgressModelFromEntity(gameProgress)

	result := r.db.First(progressModel)
	if result.Error != nil {
		return nil, result.Error
	}

	progressModel.TargetPlayerId = nextPlayerID
	r.db.Save(progressModel)

	return progressModel.ToEntity(), nil
}
