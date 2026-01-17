package model

import (
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"gorm.io/gorm"
)

// PlayerCardModel GORM 玩家手牌模型
type PlayerCardModel struct {
	gorm.Model
	Id        int `gorm:"primaryKey;auto_increment"`
	PlayerId  int
	GameId    int
	CardId    int
	Type      string
	CreatedAt time.Time   `gorm:"autoCreateTime"`
	UpdatedAt time.Time   `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt
	Card      CardModel   `gorm:"foreignKey:CardId"`
	Player    PlayerModel `gorm:"foreignKey:PlayerId"`
}

// TableName 指定表名
func (PlayerCardModel) TableName() string {
	return "player_cards"
}

// ToEntity 轉換為領域實體
func (m *PlayerCardModel) ToEntity() *entity.PlayerCard {
	return &entity.PlayerCard{
		ID:        m.Id,
		PlayerID:  m.PlayerId,
		GameID:    m.GameId,
		CardID:    m.CardId,
		Type:      m.Type,
		Card:      *m.Card.ToEntity(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// PlayerCardModelFromEntity 從領域實體轉換
func PlayerCardModelFromEntity(e *entity.PlayerCard) *PlayerCardModel {
	return &PlayerCardModel{
		Id:        e.ID,
		PlayerId:  e.PlayerID,
		GameId:    e.GameID,
		CardId:    e.CardID,
		Type:      e.Type,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
