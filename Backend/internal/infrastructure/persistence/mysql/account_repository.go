package mysql

import (
	"context"
	"strings"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/persistence/mysql/model"
	"gorm.io/gorm"
)

// AccountRepository MySQL 帳號倉儲實作
type AccountRepository struct {
	db *gorm.DB
}

// NewAccountRepository 建立帳號倉儲
func NewAccountRepository(db *gorm.DB) repository.AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// CreateAccount 建立帳號
func (r *AccountRepository) CreateAccount(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	accountModel := model.AccountModelFromEntity(account)

	err := r.db.Create(accountModel).Error
	if err != nil {
		return nil, err
	}

	return accountModel.ToEntity(), nil
}

// GetAccountById 根據 ID 查詢帳號
func (r *AccountRepository) GetAccountById(ctx context.Context, id int) (*entity.Account, error) {
	var accountModel model.AccountModel

	result := r.db.First(&accountModel, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return accountModel.ToEntity(), nil
}

// GetAccountByName 根據名稱查詢帳號
func (r *AccountRepository) GetAccountByName(ctx context.Context, name string) (*entity.Account, error) {
	var accountModel model.AccountModel

	result := r.db.Where("LOWER(name) = LOWER(?)", name).First(&accountModel)
	if result.Error != nil {
		return nil, result.Error
	}

	return accountModel.ToEntity(), nil
}

// ExistsByName 檢查名稱是否已存在（不區分大小寫）
func (r *AccountRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64

	// 使用 LOWER 進行不區分大小寫比較
	result := r.db.Model(&model.AccountModel{}).
		Where("LOWER(name) = LOWER(?)", strings.TrimSpace(name)).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}
