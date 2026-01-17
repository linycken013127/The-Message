package model

import (
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"gorm.io/gorm"
)

// CardModel GORM 卡片模型
type CardModel struct {
	gorm.Model
	Id               int `gorm:"primaryKey;auto_increment"`
	Name             string
	Color            string
	IntelligenceType int
	PlayerCards      []PlayerCardModel `gorm:"foreignKey:CardId"`
	CreatedAt        time.Time         `gorm:"autoCreateTime"`
	UpdatedAt        time.Time         `gorm:"autoCreateTime"`
	DeletedAt        gorm.DeletedAt
}

// TableName 指定表名
func (CardModel) TableName() string {
	return "cards"
}

// ToEntity 轉換為領域實體
func (m *CardModel) ToEntity() *entity.Card {
	return &entity.Card{
		ID:               m.Id,
		Name:             m.Name,
		Color:            m.Color,
		IntelligenceType: m.IntelligenceType,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

// CardModelFromEntity 從領域實體轉換
func CardModelFromEntity(e *entity.Card) *CardModel {
	return &CardModel{
		Id:               e.ID,
		Name:             e.Name,
		Color:            e.Color,
		IntelligenceType: e.IntelligenceType,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}
