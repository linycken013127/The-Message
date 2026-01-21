package repository

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// IntelligenceTransferRepository 情報傳遞儲存庫介面
type IntelligenceTransferRepository interface {
	// CreateIntelligenceTransfer 建立情報傳遞
	CreateIntelligenceTransfer(ctx context.Context, transfer *entity.IntelligenceTransfer) (*entity.IntelligenceTransfer, error)
	// GetIntelligenceTransferByID 根據 ID 取得情報傳遞
	GetIntelligenceTransferByID(ctx context.Context, id int) (*entity.IntelligenceTransfer, error)
	// GetActiveTransferByGameID 取得遊戲中正在傳遞的情報
	GetActiveTransferByGameID(ctx context.Context, gameID int) (*entity.IntelligenceTransfer, error)
	// UpdateIntelligenceTransfer 更新情報傳遞
	UpdateIntelligenceTransfer(ctx context.Context, transfer *entity.IntelligenceTransfer) error
	// DeleteIntelligenceTransfer 刪除情報傳遞
	DeleteIntelligenceTransfer(ctx context.Context, id int) error
}
