package model

import (
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"gorm.io/gorm"
)

// PlayerModel GORM 玩家模型
type PlayerModel struct {
	gorm.Model
	Id           int `gorm:"primaryKey;auto_increment"`
	Name         string
	GameId       int               `gorm:"foreignKey:GameId;references:Id"`
	IdentityCard string
	Status       string
	OrderNumber  int
	PlayerCards  []PlayerCardModel `gorm:"foreignKey:PlayerId"`
	Game         *GameModel        `gorm:"foreignKey:GameId;references:Id"`
	CreatedAt    time.Time         `gorm:"autoCreateTime"`
	UpdatedAt    time.Time         `gorm:"autoCreateTime"`
	DeletedAt    gorm.DeletedAt
}

// TableName 指定表名
func (PlayerModel) TableName() string {
	return "players"
}

// ToEntity 轉換為領域實體
func (m *PlayerModel) ToEntity() *entity.Player {
	player := &entity.Player{
		ID:           m.Id,
		Name:         m.Name,
		GameID:       m.GameId,
		IdentityCard: m.IdentityCard,
		Status:       m.Status,
		OrderNumber:  m.OrderNumber,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		PlayerCards:  make([]entity.PlayerCard, 0, len(m.PlayerCards)),
	}

	for _, pc := range m.PlayerCards {
		player.PlayerCards = append(player.PlayerCards, *pc.ToEntity())
	}

	if m.Game != nil {
		player.Game = m.Game.ToEntity()
	}

	return player
}

// PlayerModelFromEntity 從領域實體轉換
func PlayerModelFromEntity(e *entity.Player) *PlayerModel {
	model := &PlayerModel{
		Id:           e.ID,
		Name:         e.Name,
		GameId:       e.GameID,
		IdentityCard: e.IdentityCard,
		Status:       e.Status,
		OrderNumber:  e.OrderNumber,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}

	return model
}
