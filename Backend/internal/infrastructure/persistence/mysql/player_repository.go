package mysql

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql/model"
	"gorm.io/gorm"
)

// PlayerRepository MySQL 玩家倉儲實作
type PlayerRepository struct {
	db *gorm.DB
}

// NewPlayerRepository 建立玩家倉儲
func NewPlayerRepository(db *gorm.DB) repository.PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}

// CreatePlayer 建立玩家
func (r *PlayerRepository) CreatePlayer(ctx context.Context, player *entity.Player) (*entity.Player, error) {
	playerModel := model.PlayerModelFromEntity(player)

	err := r.db.Create(playerModel).Error
	if err != nil {
		return nil, err
	}

	return playerModel.ToEntity(), nil
}

// GetPlayerById 根據 ID 查詢玩家
func (r *PlayerRepository) GetPlayerById(ctx context.Context, playerID int) (*entity.Player, error) {
	var playerModel model.PlayerModel

	result := r.db.First(&playerModel, "id = ?", playerID)
	if result.Error != nil {
		return nil, result.Error
	}

	return playerModel.ToEntity(), nil
}

// GetPlayersByGameId 根據遊戲 ID 查詢所有玩家
func (r *PlayerRepository) GetPlayersByGameId(ctx context.Context, gameID int) ([]*entity.Player, error) {
	var playerModels []model.PlayerModel

	result := r.db.Find(&playerModels, "game_id = ?", gameID)
	if result.Error != nil {
		return nil, result.Error
	}

	players := make([]*entity.Player, len(playerModels))
	for i, pm := range playerModels {
		players[i] = pm.ToEntity()
	}

	return players, nil
}

// GetPlayerWithPlayerCards 查詢玩家及其手牌
func (r *PlayerRepository) GetPlayerWithPlayerCards(ctx context.Context, playerID int) (*entity.Player, error) {
	var playerModel model.PlayerModel

	if err := r.db.Preload("PlayerCards").Preload("PlayerCards.Card").First(&playerModel, playerID).Error; err != nil {
		return nil, err
	}

	return playerModel.ToEntity(), nil
}

// GetPlayerWithGame 查詢玩家及其遊戲
func (r *PlayerRepository) GetPlayerWithGame(ctx context.Context, playerID int) (*entity.Player, error) {
	var playerModel model.PlayerModel

	if err := r.db.Preload("Game").First(&playerModel, playerID).Error; err != nil {
		return nil, err
	}

	return playerModel.ToEntity(), nil
}

// GetPlayerWithGamePlayersAndPlayerCardsCard 查詢玩家完整關聯資料
func (r *PlayerRepository) GetPlayerWithGamePlayersAndPlayerCardsCard(ctx context.Context, playerID int) (*entity.Player, error) {
	var playerModel model.PlayerModel

	if err := r.db.Preload("Game.Players.PlayerCards.Card").Preload("Game.Players.Game").Preload("PlayerCards.Card").First(&playerModel, playerID).Error; err != nil {
		return nil, err
	}

	return playerModel.ToEntity(), nil
}
