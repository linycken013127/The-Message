package repository

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// AccountRepository 帳號倉儲介面
type AccountRepository interface {
	// CreateAccount 建立帳號
	CreateAccount(ctx context.Context, account *entity.Account) (*entity.Account, error)

	// GetAccountById 根據 ID 查詢帳號
	GetAccountById(ctx context.Context, id int) (*entity.Account, error)

	// GetAccountByName 根據名稱查詢帳號
	GetAccountByName(ctx context.Context, name string) (*entity.Account, error)

	// ExistsByName 檢查名稱是否已存在（不區分大小寫）
	ExistsByName(ctx context.Context, name string) (bool, error)
}
