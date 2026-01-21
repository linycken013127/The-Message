package mysql

import (
	"context"
	"errors"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql/model"
	"gorm.io/gorm"
)

// intelligenceTransferRepository 情報傳遞儲存庫實作
type intelligenceTransferRepository struct {
	db *gorm.DB
}

// NewIntelligenceTransferRepository 建立情報傳遞儲存庫
func NewIntelligenceTransferRepository(db *gorm.DB) repository.IntelligenceTransferRepository {
	return &intelligenceTransferRepository{db: db}
}

// CreateIntelligenceTransfer 建立情報傳遞
func (r *intelligenceTransferRepository) CreateIntelligenceTransfer(ctx context.Context, transfer *entity.IntelligenceTransfer) (*entity.IntelligenceTransfer, error) {
	m := model.IntelligenceTransferModelFromEntity(transfer)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetIntelligenceTransferByID 根據 ID 取得情報傳遞
func (r *intelligenceTransferRepository) GetIntelligenceTransferByID(ctx context.Context, id int) (*entity.IntelligenceTransfer, error) {
	var m model.IntelligenceTransferModel
	if err := r.db.WithContext(ctx).Preload("Card").First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToEntity(), nil
}

// GetActiveTransferByGameID 取得遊戲中正在傳遞的情報
func (r *intelligenceTransferRepository) GetActiveTransferByGameID(ctx context.Context, gameID int) (*entity.IntelligenceTransfer, error) {
	var m model.IntelligenceTransferModel
	if err := r.db.WithContext(ctx).
		Preload("Card").
		Where("game_id = ? AND status = ?", gameID, entity.IntelligenceStatusInTransit).
		First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToEntity(), nil
}

// UpdateIntelligenceTransfer 更新情報傳遞
func (r *intelligenceTransferRepository) UpdateIntelligenceTransfer(ctx context.Context, transfer *entity.IntelligenceTransfer) error {
	m := model.IntelligenceTransferModelFromEntity(transfer)
	return r.db.WithContext(ctx).Save(m).Error
}

// DeleteIntelligenceTransfer 刪除情報傳遞
func (r *intelligenceTransferRepository) DeleteIntelligenceTransfer(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.IntelligenceTransferModel{}, id).Error
}
