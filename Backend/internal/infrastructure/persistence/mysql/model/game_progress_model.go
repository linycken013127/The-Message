package model

import (
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"gorm.io/gorm"
)

// GameProgressModel GORM 遊戲進度模型
type GameProgressModel struct {
	Id             int `gorm:"primaryKey;auto_increment"`
	PlayerId       int
	GameId         int
	CardId         int
	Action         string
	TargetPlayerId int
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoCreateTime"`
	DeletedAt      gorm.DeletedAt
}

// TableName 指定表名
func (GameProgressModel) TableName() string {
	return "game_progresses"
}

// ToEntity 轉換為領域實體
func (m *GameProgressModel) ToEntity() *entity.GameProgress {
	return &entity.GameProgress{
		ID:             m.Id,
		PlayerID:       m.PlayerId,
		GameID:         m.GameId,
		CardID:         m.CardId,
		Action:         m.Action,
		TargetPlayerID: m.TargetPlayerId,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// GameProgressModelFromEntity 從領域實體轉換
func GameProgressModelFromEntity(e *entity.GameProgress) *GameProgressModel {
	return &GameProgressModel{
		Id:             e.ID,
		PlayerId:       e.PlayerID,
		GameId:         e.GameID,
		CardId:         e.CardID,
		Action:         e.Action,
		TargetPlayerId: e.TargetPlayerID,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}
