package mysql

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql/model"
	"gorm.io/gorm"
)

// CardRepository MySQL 卡片倉儲實作
type CardRepository struct {
	db *gorm.DB
}

// NewCardRepository 建立卡片倉儲
func NewCardRepository(db *gorm.DB) repository.CardRepository {
	return &CardRepository{
		db: db,
	}
}

// GetCardById 根據 ID 查詢卡片
func (r *CardRepository) GetCardById(ctx context.Context, id int) (*entity.Card, error) {
	var cardModel model.CardModel

	result := r.db.First(&cardModel, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return cardModel.ToEntity(), nil
}

// CreateCard 建立卡片
func (r *CardRepository) CreateCard(ctx context.Context, card *entity.Card) (*entity.Card, error) {
	cardModel := model.CardModelFromEntity(card)

	result := r.db.Create(cardModel)
	if result.Error != nil {
		return nil, result.Error
	}

	return cardModel.ToEntity(), nil
}

// GetCards 取得所有卡片
func (r *CardRepository) GetCards(ctx context.Context) ([]*entity.Card, error) {
	var cardModels []model.CardModel

	result := r.db.Find(&cardModels)
	if result.Error != nil {
		return nil, result.Error
	}

	cards := make([]*entity.Card, len(cardModels))
	for i, cm := range cardModels {
		cards[i] = cm.ToEntity()
	}

	return cards, nil
}

// GetPlayerCardsByPlayerId 根據玩家 ID 取得玩家手牌
func (r *CardRepository) GetPlayerCardsByPlayerId(ctx context.Context, playerID int, gameID int) ([]*entity.Card, error) {
	var playerCardModels []model.PlayerCardModel
	var cardModels []model.CardModel

	result := r.db.Find(&playerCardModels, "player_id = ? AND game_id = ?", playerID, gameID)
	if result.Error != nil {
		return nil, result.Error
	}

	var cardIDs []int
	for _, pc := range playerCardModels {
		cardIDs = append(cardIDs, pc.CardId)
	}

	result = r.db.Find(&cardModels, "id IN ?", cardIDs)
	if result.Error != nil {
		return nil, result.Error
	}

	cards := make([]*entity.Card, len(cardModels))
	for i, cm := range cardModels {
		cards[i] = cm.ToEntity()
	}

	return cards, nil
}
