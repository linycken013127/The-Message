package model

import (
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// ActionPassModel GORM 行動跳過模型
type ActionPassModel struct {
	Id        int       `gorm:"primaryKey;auto_increment"`
	GameId    int       `gorm:"column:game_id"`
	PlayerId  int       `gorm:"column:player_id"`
	Round     int       `gorm:"column:round"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName 指定表名
func (ActionPassModel) TableName() string {
	return "action_passes"
}

// ToEntity 轉換為領域實體
func (m *ActionPassModel) ToEntity() *entity.ActionPass {
	return &entity.ActionPass{
		ID:        m.Id,
		GameID:    m.GameId,
		PlayerID:  m.PlayerId,
		Round:     m.Round,
		CreatedAt: m.CreatedAt,
	}
}

// ActionPassModelFromEntity 從領域實體轉換
func ActionPassModelFromEntity(e *entity.ActionPass) *ActionPassModel {
	return &ActionPassModel{
		Id:        e.ID,
		GameId:    e.GameID,
		PlayerId:  e.PlayerID,
		Round:     e.Round,
		CreatedAt: e.CreatedAt,
	}
}
