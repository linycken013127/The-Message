package model

import (
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"gorm.io/gorm"
)

// IntelligenceTransferModel GORM 情報傳遞模型
type IntelligenceTransferModel struct {
	gorm.Model
	Id                     int       `gorm:"primaryKey;auto_increment"`
	GameId                 int       `gorm:"column:game_id"`
	CardId                 int       `gorm:"column:card_id"`
	SenderPlayerId         int       `gorm:"column:sender_player_id"`
	CurrentTargetPlayerId  int       `gorm:"column:current_target_player_id"`
	OriginalTargetPlayerId int       `gorm:"column:original_target_player_id"`
	FaceUp                 bool      `gorm:"column:face_up"`
	Status                 string    `gorm:"column:status"`
	Card                   CardModel `gorm:"foreignKey:CardId;references:Id"`
	CreatedAt              time.Time `gorm:"autoCreateTime"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime"`
	DeletedAt              gorm.DeletedAt
}

// TableName 指定表名
func (IntelligenceTransferModel) TableName() string {
	return "intelligence_transfers"
}

// ToEntity 轉換為領域實體
func (m *IntelligenceTransferModel) ToEntity() *entity.IntelligenceTransfer {
	transfer := &entity.IntelligenceTransfer{
		ID:                     m.Id,
		GameID:                 m.GameId,
		CardID:                 m.CardId,
		SenderPlayerID:         m.SenderPlayerId,
		CurrentTargetPlayerID:  m.CurrentTargetPlayerId,
		OriginalTargetPlayerID: m.OriginalTargetPlayerId,
		FaceUp:                 m.FaceUp,
		Status:                 m.Status,
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
	}

	if m.Card.Id != 0 {
		transfer.Card = m.Card.ToEntity()
	}

	return transfer
}

// IntelligenceTransferModelFromEntity 從領域實體轉換
func IntelligenceTransferModelFromEntity(e *entity.IntelligenceTransfer) *IntelligenceTransferModel {
	return &IntelligenceTransferModel{
		Id:                     e.ID,
		GameId:                 e.GameID,
		CardId:                 e.CardID,
		SenderPlayerId:         e.SenderPlayerID,
		CurrentTargetPlayerId:  e.CurrentTargetPlayerID,
		OriginalTargetPlayerId: e.OriginalTargetPlayerID,
		FaceUp:                 e.FaceUp,
		Status:                 e.Status,
		CreatedAt:              e.CreatedAt,
		UpdatedAt:              e.UpdatedAt,
	}
}
