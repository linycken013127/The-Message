package usecase

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
)

// CardUseCase 卡片用例介面
type CardUseCase interface {
	// GetCards 取得所有卡片
	GetCards(ctx context.Context) ([]*entity.Card, error)
	// GetPlayerCardsByPlayerId 取得玩家手牌
	GetPlayerCardsByPlayerId(ctx context.Context, playerID int) ([]*entity.Card, error)
}

// cardUseCase 卡片用例實作
type cardUseCase struct {
	cardRepo       repository.CardRepository
	gameRepo       repository.GameRepository
	playerRepo     repository.PlayerRepository
	playerCardRepo repository.PlayerCardRepository
}

// CardUseCaseOptions 卡片用例選項
type CardUseCaseOptions struct {
	CardRepo       repository.CardRepository
	GameRepo       repository.GameRepository
	PlayerRepo     repository.PlayerRepository
	PlayerCardRepo repository.PlayerCardRepository
}

// NewCardUseCase 建立卡片用例
func NewCardUseCase(opts *CardUseCaseOptions) CardUseCase {
	return &cardUseCase{
		cardRepo:       opts.CardRepo,
		gameRepo:       opts.GameRepo,
		playerRepo:     opts.PlayerRepo,
		playerCardRepo: opts.PlayerCardRepo,
	}
}

// GetCards 取得所有卡片
func (uc *cardUseCase) GetCards(ctx context.Context) ([]*entity.Card, error) {
	cards, err := uc.cardRepo.GetCards(ctx)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

// GetPlayerCardsByPlayerId 取得玩家手牌
func (uc *cardUseCase) GetPlayerCardsByPlayerId(ctx context.Context, playerID int) ([]*entity.Card, error) {
	player, err := uc.playerRepo.GetPlayerById(ctx, playerID)
	if err != nil {
		return nil, err
	}

	game, err := uc.gameRepo.GetGameWithPlayers(ctx, player.GameID)
	if err != nil {
		return nil, err
	}

	cards, err := uc.cardRepo.GetPlayerCardsByPlayerId(ctx, player.ID, game.ID)
	if err != nil {
		return nil, err
	}

	return cards, nil
}
