package model

import (
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"gorm.io/gorm"
)

// DeckModel GORM 牌組模型
type DeckModel struct {
	gorm.Model
	Id        int `gorm:"primaryKey;auto_increment"`
	GameId    int
	CardId    int
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt
}

// TableName 指定表名
func (DeckModel) TableName() string {
	return "decks"
}

// ToEntity 轉換為領域實體
func (m *DeckModel) ToEntity() *entity.Deck {
	return &entity.Deck{
		ID:        m.Id,
		GameID:    m.GameId,
		CardID:    m.CardId,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// DeckModelFromEntity 從領域實體轉換
func DeckModelFromEntity(e *entity.Deck) *DeckModel {
	return &DeckModel{
		Id:        e.ID,
		GameId:    e.GameID,
		CardId:    e.CardID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
