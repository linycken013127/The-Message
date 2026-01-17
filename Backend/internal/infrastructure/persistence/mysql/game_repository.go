package mysql

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql/model"
	"gorm.io/gorm"
)

// GameRepository MySQL 遊戲倉儲實作
type GameRepository struct {
	db *gorm.DB
}

// NewGameRepository 建立遊戲倉儲
func NewGameRepository(db *gorm.DB) repository.GameRepository {
	return &GameRepository{
		db: db,
	}
}

// GetGameById 根據 ID 查詢遊戲
func (r *GameRepository) GetGameById(ctx context.Context, id int) (*entity.Game, error) {
	var gameModel model.GameModel

	result := r.db.First(&gameModel, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return gameModel.ToEntity(), nil
}

// CreateGame 建立遊戲
func (r *GameRepository) CreateGame(ctx context.Context, game *entity.Game) (*entity.Game, error) {
	gameModel := model.GameModelFromEntity(game)

	result := r.db.Create(gameModel)
	if result.Error != nil {
		return nil, result.Error
	}

	return gameModel.ToEntity(), nil
}

// DeleteGame 刪除遊戲
func (r *GameRepository) DeleteGame(ctx context.Context, id int) error {
	result := r.db.Delete(&model.GameModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetGameWithPlayers 查詢遊戲及其玩家
func (r *GameRepository) GetGameWithPlayers(ctx context.Context, id int) (*entity.Game, error) {
	var gameModel model.GameModel

	if err := r.db.Preload("Players").First(&gameModel, id).Error; err != nil {
		return nil, err
	}

	return gameModel.ToEntity(), nil
}

// UpdateGame 更新遊戲
func (r *GameRepository) UpdateGame(ctx context.Context, game *entity.Game) error {
	gameModel := model.GameModelFromEntity(game)

	result := r.db.Save(gameModel)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
