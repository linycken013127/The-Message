package mysql

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql/model"
	"gorm.io/gorm"
)

// DeckRepository MySQL 牌組倉儲實作
type DeckRepository struct {
	db *gorm.DB
}

// NewDeckRepository 建立牌組倉儲
func NewDeckRepository(db *gorm.DB) repository.DeckRepository {
	return &DeckRepository{
		db: db,
	}
}

// CreateDeck 建立牌組
func (r *DeckRepository) CreateDeck(ctx context.Context, deck *entity.Deck) (*entity.Deck, error) {
	deckModel := model.DeckModelFromEntity(deck)

	result := r.db.Create(deckModel)
	if result.Error != nil {
		return nil, result.Error
	}

	return deckModel.ToEntity(), nil
}

// GetDecksByGameId 根據遊戲 ID 查詢牌組
func (r *DeckRepository) GetDecksByGameId(ctx context.Context, gameID int) ([]*entity.Deck, error) {
	var deckModels []model.DeckModel

	result := r.db.Find(&deckModels, "game_id = ?", gameID)
	if result.Error != nil {
		return nil, result.Error
	}

	decks := make([]*entity.Deck, len(deckModels))
	for i, dm := range deckModels {
		decks[i] = dm.ToEntity()
	}

	return decks, nil
}

// DeleteDeck 刪除牌組
func (r *DeckRepository) DeleteDeck(ctx context.Context, id int) error {
	result := r.db.Delete(&model.DeckModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
