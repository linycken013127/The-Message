package model

import (
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"gorm.io/gorm"
)

// GameModel GORM 遊戲模型
type GameModel struct {
	gorm.Model
	Id              int `gorm:"primaryKey;auto_increment"`
	Token           string
	Status          string
	Phase           string
	CurrentPlayerId int
	HostAccountId   int
	MaxPlayers      int
	CurrentPlayers  int
	Winner          string
	Players         []PlayerModel     `gorm:"foreignKey:GameId"`
	GamePlayers     []GamePlayerModel `gorm:"foreignKey:GameId"`
	CreatedAt       time.Time         `gorm:"autoCreateTime"`
	UpdatedAt       time.Time         `gorm:"autoCreateTime"`
	DeletedAt       gorm.DeletedAt
}

// TableName 指定表名
func (GameModel) TableName() string {
	return "games"
}

// ToEntity 轉換為領域實體
func (m *GameModel) ToEntity() *entity.Game {
	game := &entity.Game{
		ID:              m.Id,
		Token:           m.Token,
		Status:          m.Status,
		Phase:           m.Phase,
		CurrentPlayerID: m.CurrentPlayerId,
		HostAccountID:   m.HostAccountId,
		MaxPlayers:      m.MaxPlayers,
		CurrentPlayers:  m.CurrentPlayers,
		Winner:          m.Winner,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		Players:         make([]entity.Player, 0, len(m.Players)),
		GamePlayers:     make([]entity.GamePlayer, 0, len(m.GamePlayers)),
	}

	for _, pm := range m.Players {
		player := pm.ToEntity()
		player.Game = game
		game.Players = append(game.Players, *player)
	}

	for _, gpm := range m.GamePlayers {
		game.GamePlayers = append(game.GamePlayers, *gpm.ToEntity())
	}

	return game
}

// GameModelFromEntity 從領域實體轉換
func GameModelFromEntity(e *entity.Game) *GameModel {
	model := &GameModel{
		Id:              e.ID,
		Token:           e.Token,
		Status:          e.Status,
		Phase:           e.Phase,
		CurrentPlayerId: e.CurrentPlayerID,
		HostAccountId:   e.HostAccountID,
		MaxPlayers:      e.MaxPlayers,
		CurrentPlayers:  e.CurrentPlayers,
		Winner:          e.Winner,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}

	return model
}
