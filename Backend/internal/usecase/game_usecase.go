package usecase

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"math/rand"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
)

// GameUseCase 遊戲用例介面
type GameUseCase interface {
	// StartGame 開始遊戲（初始化遊戲、玩家、牌組、抽牌）
	StartGame(ctx context.Context, req CreateGameRequest) (*entity.Game, error)
	// InitGame 初始化遊戲
	InitGame(ctx context.Context) (*entity.Game, error)
	// InitDeck 初始化牌組
	InitDeck(ctx context.Context, game *entity.Game) error
	// DrawCard 抽牌
	DrawCard(ctx context.Context, game *entity.Game, player *entity.Player, drawCards []*entity.Deck, count int) error
	// DrawCardsForAllPlayers 為所有玩家抽牌
	DrawCardsForAllPlayers(ctx context.Context, game *entity.Game) error
	// GetGameById 根據 ID 取得遊戲
	GetGameById(ctx context.Context, id int) (*entity.Game, error)
	// CreateGame 建立遊戲
	CreateGame(ctx context.Context, game *entity.Game) (*entity.Game, error)
	// DeleteGame 刪除遊戲
	DeleteGame(ctx context.Context, id int) error
	// UpdateCurrentPlayer 更新當前玩家
	UpdateCurrentPlayer(ctx context.Context, game *entity.Game, playerID int)
	// NextPlayer 切換到下一位玩家
	NextPlayer(ctx context.Context, player *entity.Player) (*entity.Game, error)
	// UpdateStatus 更新遊戲狀態
	UpdateStatus(ctx context.Context, game *entity.Game, status string)
}

// gameUseCase 遊戲用例實作
type gameUseCase struct {
	gameRepo       repository.GameRepository
	playerUseCase  PlayerUseCase
	cardUseCase    CardUseCase
	deckUseCase    DeckUseCase
	eventPublisher repository.EventPublisher
}

// GameUseCaseOptions 遊戲用例選項
type GameUseCaseOptions struct {
	GameRepo       repository.GameRepository
	PlayerUseCase  PlayerUseCase
	CardUseCase    CardUseCase
	DeckUseCase    DeckUseCase
	EventPublisher repository.EventPublisher
}

// NewGameUseCase 建立遊戲用例
func NewGameUseCase(opts *GameUseCaseOptions) GameUseCase {
	return &gameUseCase{
		gameRepo:       opts.GameRepo,
		playerUseCase:  opts.PlayerUseCase,
		cardUseCase:    opts.CardUseCase,
		deckUseCase:    opts.DeckUseCase,
		eventPublisher: opts.EventPublisher,
	}
}

// StartGame 開始遊戲（初始化遊戲、玩家、牌組、抽牌）
func (uc *gameUseCase) StartGame(ctx context.Context, req CreateGameRequest) (*entity.Game, error) {
	// 1. 初始化遊戲
	game, err := uc.InitGame(ctx)
	if err != nil {
		return nil, err
	}

	// 2. 初始化玩家
	if err := uc.playerUseCase.InitPlayers(ctx, game, req); err != nil {
		return nil, err
	}

	// 3. 重新取得遊戲（包含玩家資料）
	game, err = uc.GetGameById(ctx, game.ID)
	if err != nil {
		return nil, err
	}

	// 4. 設定當前玩家和遊戲狀態
	uc.UpdateCurrentPlayer(ctx, game, game.Players[0].ID)
	uc.UpdateStatus(ctx, game, entity.GameStatusActionCardStage)

	// 5. 初始化牌組
	if err := uc.InitDeck(ctx, game); err != nil {
		return nil, err
	}

	// 6. 為所有玩家抽牌
	if err := uc.DrawCardsForAllPlayers(ctx, game); err != nil {
		return nil, err
	}

	// 7. 回傳最新的遊戲資料
	game, err = uc.GetGameById(ctx, game.ID)
	if err != nil {
		return nil, err
	}

	// 8. 發送遊戲開始事件
	if uc.eventPublisher != nil {
		uc.eventPublisher.PublishGameStarted(repository.GameEvent{
			GameID:     game.ID,
			Message:    "Game started",
			NextPlayer: game.Players[0].ID,
		})
	}

	return game, nil
}

// InitGame 初始化遊戲
func (uc *gameUseCase) InitGame(ctx context.Context) (*entity.Game, error) {
	token, err := generateSecureToken(256)
	if err != nil {
		return nil, err
	}

	game, err := uc.CreateGame(ctx, entity.NewGame(token, entity.GameStatusStart))
	if err != nil {
		return nil, err
	}

	return game, nil
}

// InitDeck 初始化牌組
func (uc *gameUseCase) InitDeck(ctx context.Context, game *entity.Game) error {
	return uc.deckUseCase.InitDeck(ctx, game)
}

// DrawCard 抽牌
func (uc *gameUseCase) DrawCard(ctx context.Context, game *entity.Game, player *entity.Player, drawCards []*entity.Deck, count int) error {
	for i := 0; i < count; i++ {
		card := entity.NewPlayerCard(player.ID, game.ID, drawCards[i].CardID, entity.PlayerCardTypeHand)
		err := uc.playerUseCase.CreatePlayerCard(ctx, card)
		if err != nil {
			return err
		}
		err = uc.deckUseCase.DeleteDeckFromGame(ctx, drawCards[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}

// DrawCardsForAllPlayers 為所有玩家抽牌
func (uc *gameUseCase) DrawCardsForAllPlayers(ctx context.Context, game *entity.Game) error {
	players, err := uc.playerUseCase.GetPlayersByGameId(ctx, game.ID)
	if err != nil {
		return err
	}
	for _, player := range players {
		drawCards, _ := uc.deckUseCase.GetDecksByGameId(ctx, game.ID)
		err := uc.DrawCard(ctx, game, player, drawCards, 3)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetGameById 根據 ID 取得遊戲
func (uc *gameUseCase) GetGameById(ctx context.Context, id int) (*entity.Game, error) {
	game, err := uc.gameRepo.GetGameWithPlayers(ctx, id)
	if err != nil {
		return nil, err
	}
	return game, nil
}

// CreateGame 建立遊戲
func (uc *gameUseCase) CreateGame(ctx context.Context, game *entity.Game) (*entity.Game, error) {
	game, err := uc.gameRepo.CreateGame(ctx, game)
	if err != nil {
		return nil, err
	}
	return game, nil
}

// DeleteGame 刪除遊戲
func (uc *gameUseCase) DeleteGame(ctx context.Context, id int) error {
	return uc.gameRepo.DeleteGame(ctx, id)
}

// UpdateCurrentPlayer 更新當前玩家
func (uc *gameUseCase) UpdateCurrentPlayer(ctx context.Context, game *entity.Game, playerID int) {
	game.CurrentPlayerID = playerID
	err := uc.gameRepo.UpdateGame(ctx, game)
	if err != nil {
		panic(err)
	}
}

// NextPlayer 切換到下一位玩家
func (uc *gameUseCase) NextPlayer(ctx context.Context, player *entity.Player) (*entity.Game, error) {
	players := player.Game.Players
	currentPlayerID := player.ID

	var currentPlayerIndex int
	for index, gPlayer := range players {
		if gPlayer.ID == currentPlayerID {
			currentPlayerIndex = index
			break
		}
	}

	if currentPlayerIndex+1 >= len(players) {
		player.Game.CurrentPlayerID = players[0].ID
		player.Game.Status = entity.GameStatusTransmitIntelligence
	} else {
		player.Game.CurrentPlayerID = players[currentPlayerIndex+1].ID
	}
	return player.Game, nil
}

// UpdateStatus 更新遊戲狀態
func (uc *gameUseCase) UpdateStatus(ctx context.Context, game *entity.Game, status string) {
	game.Status = status
	err := uc.gameRepo.UpdateGame(ctx, game)
	if err != nil {
		panic(err)
	}
}

// generateSecureToken 產生安全 Token
func generateSecureToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := cryptorand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ShuffleIdentityCards 洗身份卡
func ShuffleIdentityCards(cards []string) []string {
	shuffledCards := make([]string, len(cards))
	perm := rand.Perm(len(cards))
	for i, j := range perm {
		shuffledCards[i] = cards[j]
	}
	return shuffledCards
}
