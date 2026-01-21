package usecase

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
)

// DeathUseCase 死亡條件用例介面
type DeathUseCase interface {
	// CheckAndHandleDeath 檢查並處理玩家死亡
	CheckAndHandleDeath(ctx context.Context, playerID int, gameID int) (bool, error)
	// KillPlayer 殺死玩家
	KillPlayer(ctx context.Context, playerID int) error
	// CheckAllPlayersDead 檢查是否所有玩家都死亡
	CheckAllPlayersDead(ctx context.Context, gameID int) (bool, error)
	// CountPlayerBlackIntelligence 計算玩家黑色情報數量
	CountPlayerBlackIntelligence(ctx context.Context, playerID int) (int, error)
	// EndGameAsDraw 以平局結束遊戲
	EndGameAsDraw(ctx context.Context, gameID int) error
}

// deathUseCase 死亡條件用例實作
type deathUseCase struct {
	playerRepo     repository.PlayerRepository
	playerCardRepo repository.PlayerCardRepository
	gameRepo       repository.GameRepository
	cardRepo       repository.CardRepository
}

// DeathUseCaseOptions 死亡條件用例選項
type DeathUseCaseOptions struct {
	PlayerRepo     repository.PlayerRepository
	PlayerCardRepo repository.PlayerCardRepository
	GameRepo       repository.GameRepository
	CardRepo       repository.CardRepository
}

// NewDeathUseCase 建立死亡條件用例
func NewDeathUseCase(opts *DeathUseCaseOptions) DeathUseCase {
	return &deathUseCase{
		playerRepo:     opts.PlayerRepo,
		playerCardRepo: opts.PlayerCardRepo,
		gameRepo:       opts.GameRepo,
		cardRepo:       opts.CardRepo,
	}
}

// CheckAndHandleDeath 檢查並處理玩家死亡
// 如果玩家黑色情報 >= 3，則殺死玩家並檢查是否所有玩家都死亡
func (uc *deathUseCase) CheckAndHandleDeath(ctx context.Context, playerID int, gameID int) (bool, error) {
	// 計算黑色情報數量
	blackCount, err := uc.CountPlayerBlackIntelligence(ctx, playerID)
	if err != nil {
		return false, err
	}

	// 黑色情報 >= 3 則死亡
	if blackCount >= 3 {
		err = uc.KillPlayer(ctx, playerID)
		if err != nil {
			return false, err
		}

		// 檢查是否所有玩家都死亡
		allDead, err := uc.CheckAllPlayersDead(ctx, gameID)
		if err != nil {
			return true, err
		}

		if allDead {
			err = uc.EndGameAsDraw(ctx, gameID)
			if err != nil {
				return true, err
			}
		}

		return true, nil
	}

	return false, nil
}

// KillPlayer 殺死玩家
func (uc *deathUseCase) KillPlayer(ctx context.Context, playerID int) error {
	player, err := uc.playerRepo.GetPlayerById(ctx, playerID)
	if err != nil {
		return err
	}

	player.Status = entity.PlayerStatusDead
	return uc.playerRepo.UpdatePlayer(ctx, player)
}

// CheckAllPlayersDead 檢查是否所有玩家都死亡
func (uc *deathUseCase) CheckAllPlayersDead(ctx context.Context, gameID int) (bool, error) {
	players, err := uc.playerRepo.GetPlayersByGameId(ctx, gameID)
	if err != nil {
		return false, err
	}

	for _, player := range players {
		if player.Status != entity.PlayerStatusDead {
			return false, nil
		}
	}

	return true, nil
}

// CountPlayerBlackIntelligence 計算玩家黑色情報數量
func (uc *deathUseCase) CountPlayerBlackIntelligence(ctx context.Context, playerID int) (int, error) {
	player, err := uc.playerRepo.GetPlayerWithPlayerCards(ctx, playerID)
	if err != nil {
		return 0, err
	}

	blackCount := 0
	for _, pc := range player.PlayerCards {
		if pc.Type == entity.PlayerCardTypeIntelligence && pc.Card.Color == entity.CardColorBlack {
			blackCount++
		}
	}

	return blackCount, nil
}

// EndGameAsDraw 以平局結束遊戲
func (uc *deathUseCase) EndGameAsDraw(ctx context.Context, gameID int) error {
	game, err := uc.gameRepo.GetGameById(ctx, gameID)
	if err != nil {
		return err
	}

	game.Status = entity.GameRoomStatusEnded
	game.Winner = "" // 空字串表示平局

	return uc.gameRepo.UpdateGame(ctx, game)
}
