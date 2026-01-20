package model

import (
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"gorm.io/gorm"
)

// AccountModel GORM 帳號模型
type AccountModel struct {
	gorm.Model
	Id        int `gorm:"primaryKey;auto_increment"`
	Name      string
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName 指定表名
func (AccountModel) TableName() string {
	return "accounts"
}

// ToEntity 轉換為領域實體
func (m *AccountModel) ToEntity() *entity.Account {
	return &entity.Account{
		ID:        m.Id,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// AccountModelFromEntity 從領域實體轉換
func AccountModelFromEntity(e *entity.Account) *AccountModel {
	return &AccountModel{
		Id:        e.ID,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
