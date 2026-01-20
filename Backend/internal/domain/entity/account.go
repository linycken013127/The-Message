package entity

import (
	"errors"
	"strings"
	"time"
)

// Account 帳號領域實體
// 代表一個獨立於遊戲的註冊玩家
type Account struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AccountNameMaxLength 玩家名稱最大長度
const AccountNameMaxLength = 10

// NewAccount 建立新帳號（工廠方法）
func NewAccount(name string) (*Account, error) {
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Validate name
	if err := ValidateAccountName(name); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Account{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ValidateAccountName 驗證帳號名稱
func ValidateAccountName(name string) error {
	if name == "" {
		return errors.New("玩家名稱不可為空")
	}

	// 計算字元數（支援中文等 Unicode 字元）
	runeCount := len([]rune(name))
	if runeCount > AccountNameMaxLength {
		return errors.New("玩家名稱最大長度為 10 字元")
	}

	return nil
}

// 帳號相關錯誤
var (
	ErrAccountNameEmpty     = errors.New("玩家名稱不可為空")
	ErrAccountNameTooLong   = errors.New("玩家名稱最大長度為 10 字元")
	ErrAccountNameDuplicate = errors.New("玩家名稱已存在")
)
