package entity

import "time"

// Card 卡片領域實體
type Card struct {
	ID               int
	Name             string
	Color            string
	IntelligenceType int
	PlayerCards      []PlayerCard
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewCard 建立新卡片（工廠方法）
func NewCard(name string, color string, intelligenceType int) *Card {
	now := time.Now()
	return &Card{
		Name:             name,
		Color:            color,
		IntelligenceType: intelligenceType,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// 卡片顏色常數
const (
	CardColorRed   = "紅"
	CardColorBlue  = "藍"
	CardColorBlack = "黑"
)

// 情報類型常數
const (
	IntelligenceTypeSecretTelegram = 1 // 密電
	IntelligenceTypeDirect         = 2 // 直達
	IntelligenceTypeDocument       = 3 // 文件
)

// 卡片名稱常數
const (
	CardNameLockOn      = "鎖定"
	CardNameLureAway    = "調虎離山"
	CardNameProbe       = "試探"
	CardNameIntercept   = "截獲"
	CardNameDecipher    = "破譯"
	CardNameDiversion   = "退回"
	CardNameBurn        = "燒毀"
	CardNameBlurOfTruth = "真偽莫辯"
	CardNameSeeThrough  = "識破"
)

// ToIntelligenceType 根據卡片名稱取得情報類型
func ToIntelligenceType(cardName string) int {
	switch cardName {
	case CardNameLockOn, CardNameLureAway, CardNameProbe, CardNameDecipher:
		return IntelligenceTypeSecretTelegram
	case CardNameIntercept, CardNameBurn, CardNameSeeThrough:
		return IntelligenceTypeDirect
	case CardNameDiversion, CardNameBlurOfTruth:
		return IntelligenceTypeDocument
	default:
		return 0
	}
}

// IntelligenceTypeToString 情報類型轉字串
func IntelligenceTypeToString(intelligenceType int) string {
	switch intelligenceType {
	case IntelligenceTypeSecretTelegram:
		return "密電"
	case IntelligenceTypeDirect:
		return "直達"
	case IntelligenceTypeDocument:
		return "文件"
	default:
		return ""
	}
}
