package entity

import "time"

// Player 玩家領域實體
type Player struct {
	ID           int
	Name         string
	GameID       int
	IdentityCard string
	Status       string
	OrderNumber  int
	PlayerCards  []PlayerCard
	Game         *Game
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewPlayer 建立新玩家（工廠方法）
func NewPlayer(name string, gameID int, identityCard string, orderNumber int) *Player {
	now := time.Now()
	return &Player{
		Name:         name,
		GameID:       gameID,
		IdentityCard: identityCard,
		Status:       PlayerStatusAlive,
		OrderNumber:  orderNumber,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsAlive 檢查玩家是否存活
func (p *Player) IsAlive() bool {
	return p.Status == PlayerStatusAlive
}

// IsDead 檢查玩家是否死亡
func (p *Player) IsDead() bool {
	return p.Status == PlayerStatusDead
}

// 玩家狀態常數
const (
	PlayerStatusAlive = "生存"
	PlayerStatusDead  = "死亡"
)

// 身份卡常數
const (
	IdentityUndercoverFront = "潛伏戰線"
	IdentityMilitaryAgency  = "軍情處"
	IdentityBystander       = "打醬油"
)
