package entity

import (
	"errors"
	"time"
)

// Game 遊戲領域實體
// 純 Go 結構體，不依賴任何外部框架
type Game struct {
	ID              int
	Token           string
	Status          string
	CurrentPlayerID int
	HostAccountID   int
	MaxPlayers      int
	CurrentPlayers  int
	Players         []Player
	GamePlayers     []GamePlayer
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewGame 建立新遊戲（工廠方法）
func NewGame(token string, status string) *Game {
	now := time.Now()
	return &Game{
		Token:     token,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewGameRoom 建立新遊戲房（工廠方法）
func NewGameRoom(hostAccountID int, maxPlayers int) *Game {
	now := time.Now()
	if maxPlayers <= 0 {
		maxPlayers = GameMaxPlayers
	}
	if maxPlayers > GameMaxPlayers {
		maxPlayers = GameMaxPlayers
	}
	if maxPlayers < GameMinPlayers {
		maxPlayers = GameMinPlayers
	}

	return &Game{
		Token:          "",
		Status:         GameRoomStatusWaiting,
		HostAccountID:  hostAccountID,
		MaxPlayers:     maxPlayers,
		CurrentPlayers: 1, // 房主自動加入
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// CanJoin 檢查是否可以加入
func (g *Game) CanJoin() error {
	if g.Status != GameRoomStatusWaiting {
		return ErrGameAlreadyStarted
	}
	if g.CurrentPlayers >= g.MaxPlayers {
		return ErrGameRoomFull
	}
	return nil
}

// CanStart 檢查是否可以開始遊戲
func (g *Game) CanStart(accountID int) error {
	if g.HostAccountID != accountID {
		return ErrNotGameHost
	}
	if g.CurrentPlayers < GameMinPlayers {
		return ErrNotEnoughPlayers
	}
	if g.Status != GameRoomStatusWaiting {
		return ErrGameAlreadyStarted
	}
	return nil
}

// IsWaiting 檢查遊戲是否在等待中
func (g *Game) IsWaiting() bool {
	return g.Status == GameRoomStatusWaiting
}

// IsPlaying 檢查遊戲是否在進行中
func (g *Game) IsPlaying() bool {
	return g.Status == GameRoomStatusPlaying
}

// SetCurrentPlayer 設定當前玩家
func (g *Game) SetCurrentPlayer(playerID int) {
	g.CurrentPlayerID = playerID
	g.UpdatedAt = time.Now()
}

// SetStatus 設定遊戲狀態
func (g *Game) SetStatus(status string) {
	g.Status = status
	g.UpdatedAt = time.Now()
}

// IsEnded 檢查遊戲是否結束
func (g *Game) IsEnded() bool {
	return g.Status == GameStatusEnd
}

// 遊戲狀態常數（舊版，遊戲進行中的狀態）
const (
	GameStatusStart                = "開始遊戲"
	GameStatusActionCardStage      = "功能牌階段"
	GameStatusTransmitIntelligence = "情報牌階段"
	GameStatusEnd                  = "結束遊戲"
)

// 遊戲房狀態常數
const (
	GameRoomStatusWaiting = "WAITING" // 等待中
	GameRoomStatusPlaying = "PLAYING" // 進行中
	GameRoomStatusEnded   = "ENDED"   // 已結束
)

// 遊戲房人數限制
const (
	GameMinPlayers = 3
	GameMaxPlayers = 9
)

// 遊戲房相關錯誤
var (
	ErrGameNotFound       = errors.New("遊戲房不存在")
	ErrGameAlreadyStarted = errors.New("遊戲已開始")
	ErrGameRoomFull       = errors.New("遊戲房已滿")
	ErrNotGameHost        = errors.New("只有房主可以開始遊戲")
	ErrNotEnoughPlayers   = errors.New("人數不足，至少需要 3 人")
	ErrAlreadyInGame      = errors.New("你已經在這個遊戲房中")
)
