package mysql

import (
	"context"
	"errors"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql/model"
	"gorm.io/gorm"
)

// PlayerCardRepository MySQL 玩家手牌倉儲實作
type PlayerCardRepository struct {
	db *gorm.DB
}

// NewPlayerCardRepository 建立玩家手牌倉儲
func NewPlayerCardRepository(db *gorm.DB) repository.PlayerCardRepository {
	return &PlayerCardRepository{
		db: db,
	}
}

// CreatePlayerCard 建立玩家手牌
func (r *PlayerCardRepository) CreatePlayerCard(ctx context.Context, card *entity.PlayerCard) (*entity.PlayerCard, error) {
	cardModel := model.PlayerCardModelFromEntity(card)

	result := r.db.Create(cardModel)
	if result.Error != nil {
		return nil, result.Error
	}

	return cardModel.ToEntity(), nil
}

// DeletePlayerCard 刪除玩家手牌
func (r *PlayerCardRepository) DeletePlayerCard(ctx context.Context, id int) error {
	result := r.db.Delete(&model.PlayerCardModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeletePlayerCardByPlayerIdAndCardId 根據玩家 ID 和卡片 ID 刪除手牌
func (r *PlayerCardRepository) DeletePlayerCardByPlayerIdAndCardId(ctx context.Context, playerID int, gameID int, cardID int) (bool, error) {
	result := r.db.Delete(&model.PlayerCardModel{}, "player_id = ? AND game_id = ? AND card_id = ?", playerID, gameID, cardID)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// ExistPlayerCardByPlayerIdAndCardId 檢查玩家手牌是否存在
func (r *PlayerCardRepository) ExistPlayerCardByPlayerIdAndCardId(ctx context.Context, playerID int, gameID int, cardID int) (bool, error) {
	var cardModel model.PlayerCardModel

	result := r.db.First(&cardModel, "player_id = ? AND game_id = ? AND card_id = ?", playerID, gameID, cardID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

// GetPlayerCardById 根據 ID 查詢玩家手牌
func (r *PlayerCardRepository) GetPlayerCardById(ctx context.Context, id int) (*entity.PlayerCard, error) {
	var cardModel model.PlayerCardModel

	result := r.db.First(&cardModel, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return cardModel.ToEntity(), nil
}

// GetPlayerCards 取得玩家手牌
func (r *PlayerCardRepository) GetPlayerCards(ctx context.Context, playerCard *entity.PlayerCard) (*[]entity.PlayerCard, error) {
	var playerCardModels []model.PlayerCardModel
	queryModel := model.PlayerCardModelFromEntity(playerCard)

	result := r.db.Model(&model.PlayerCardModel{}).Preload("Card").Where(queryModel).Find(&playerCardModels)
	if result.Error != nil {
		return nil, result.Error
	}

	playerCards := make([]entity.PlayerCard, len(playerCardModels))
	for i, pcm := range playerCardModels {
		playerCards[i] = *pcm.ToEntity()
	}

	return &playerCards, nil
}

// GetPlayerCardsByGameId 根據遊戲 ID 查詢玩家手牌
func (r *PlayerCardRepository) GetPlayerCardsByGameId(ctx context.Context, gameID int) ([]*entity.PlayerCard, error) {
	var playerCardModels []model.PlayerCardModel

	result := r.db.Find(&playerCardModels, "game_id = ?", gameID)
	if result.Error != nil {
		return nil, result.Error
	}

	playerCards := make([]*entity.PlayerCard, len(playerCardModels))
	for i, pcm := range playerCardModels {
		playerCards[i] = pcm.ToEntity()
	}

	return playerCards, nil
}
