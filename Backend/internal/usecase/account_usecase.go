package usecase

import (
	"context"
	"strings"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
)

// AccountUseCase 帳號用例介面
type AccountUseCase interface {
	// RegisterAccount 註冊帳號
	RegisterAccount(ctx context.Context, name string) (*entity.Account, error)

	// GetAccountById 根據 ID 取得帳號
	GetAccountById(ctx context.Context, id int) (*entity.Account, error)

	// GetAccountByName 根據名稱取得帳號
	GetAccountByName(ctx context.Context, name string) (*entity.Account, error)
}

// accountUseCase 帳號用例實作
type accountUseCase struct {
	accountRepo repository.AccountRepository
}

// AccountUseCaseOptions 帳號用例選項
type AccountUseCaseOptions struct {
	AccountRepo repository.AccountRepository
}

// NewAccountUseCase 建立帳號用例
func NewAccountUseCase(opts *AccountUseCaseOptions) AccountUseCase {
	return &accountUseCase{
		accountRepo: opts.AccountRepo,
	}
}

// RegisterAccount 註冊帳號
func (uc *accountUseCase) RegisterAccount(ctx context.Context, name string) (*entity.Account, error) {
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Validate name using domain logic
	if err := entity.ValidateAccountName(name); err != nil {
		return nil, err
	}

	// Check if name already exists (case insensitive)
	exists, err := uc.accountRepo.ExistsByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entity.ErrAccountNameDuplicate
	}

	// Create account using factory method
	account, err := entity.NewAccount(name)
	if err != nil {
		return nil, err
	}

	// Save to repository
	return uc.accountRepo.CreateAccount(ctx, account)
}

// GetAccountById 根據 ID 取得帳號
func (uc *accountUseCase) GetAccountById(ctx context.Context, id int) (*entity.Account, error) {
	return uc.accountRepo.GetAccountById(ctx, id)
}

// GetAccountByName 根據名稱取得帳號
func (uc *accountUseCase) GetAccountByName(ctx context.Context, name string) (*entity.Account, error) {
	return uc.accountRepo.GetAccountByName(ctx, name)
}
