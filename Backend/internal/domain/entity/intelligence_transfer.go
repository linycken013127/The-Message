package entity

import (
	"errors"
	"time"
)

// IntelligenceTransfer 情報傳遞領域實體
type IntelligenceTransfer struct {
	ID                     int
	GameID                 int
	CardID                 int
	SenderPlayerID         int
	CurrentTargetPlayerID  int
	OriginalTargetPlayerID int
	FaceUp                 bool
	Status                 string
	Card                   *Card
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

// NewIntelligenceTransfer 建立新情報傳遞（工廠方法）
func NewIntelligenceTransfer(
	gameID int,
	cardID int,
	senderPlayerID int,
	currentTargetPlayerID int,
	originalTargetPlayerID int,
	faceUp bool,
) *IntelligenceTransfer {
	now := time.Now()
	return &IntelligenceTransfer{
		GameID:                 gameID,
		CardID:                 cardID,
		SenderPlayerID:         senderPlayerID,
		CurrentTargetPlayerID:  currentTargetPlayerID,
		OriginalTargetPlayerID: originalTargetPlayerID,
		FaceUp:                 faceUp,
		Status:                 IntelligenceStatusInTransit,
		CreatedAt:              now,
		UpdatedAt:              now,
	}
}

// 情報傳遞狀態常數
const (
	IntelligenceStatusInTransit = "IN_TRANSIT" // 傳遞中
	IntelligenceStatusCompleted = "COMPLETED"  // 已完成
)

// 情報階段相關錯誤
var (
	ErrCannotTargetSelf       = errors.New("直達情報不能指定自己")
	ErrCannotTargetDeadPlayer = errors.New("直達情報不能指定已死亡的玩家")
	ErrDirectCardNeedsTarget  = errors.New("直達情報需要指定目標玩家")
	ErrNoIntelligenceInTransit = errors.New("目前沒有情報傳遞中")
	ErrNotIntelligenceTarget  = errors.New("你不是當前情報的接收者")
)

// IsInTransit 檢查情報是否在傳遞中
func (t *IntelligenceTransfer) IsInTransit() bool {
	return t.Status == IntelligenceStatusInTransit
}

// Complete 完成情報傳遞
func (t *IntelligenceTransfer) Complete() {
	t.Status = IntelligenceStatusCompleted
	t.UpdatedAt = time.Now()
}

// UpdateTarget 更新目標玩家
func (t *IntelligenceTransfer) UpdateTarget(newTargetID int) {
	t.CurrentTargetPlayerID = newTargetID
	t.UpdatedAt = time.Now()
}
