package model

import (
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// GamePlayerModel GORM 遊戲玩家關聯模型
type GamePlayerModel struct {
	Id        int `gorm:"primaryKey;auto_increment"`
	GameId    int
	AccountId int
	JoinOrder int
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName 指定表名
func (GamePlayerModel) TableName() string {
	return "game_players"
}

// ToEntity 轉換為領域實體
func (m *GamePlayerModel) ToEntity() *entity.GamePlayer {
	return &entity.GamePlayer{
		ID:        m.Id,
		GameID:    m.GameId,
		AccountID: m.AccountId,
		JoinOrder: m.JoinOrder,
		CreatedAt: m.CreatedAt,
	}
}

// GamePlayerModelFromEntity 從領域實體轉換
func GamePlayerModelFromEntity(e *entity.GamePlayer) *GamePlayerModel {
	return &GamePlayerModel{
		Id:        e.ID,
		GameId:    e.GameID,
		AccountId: e.AccountID,
		JoinOrder: e.JoinOrder,
		CreatedAt: e.CreatedAt,
	}
}
