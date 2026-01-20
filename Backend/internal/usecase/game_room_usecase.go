package usecase

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
)

// GameRoomUseCase 遊戲房用例介面
type GameRoomUseCase interface {
	// CreateGameRoom 建立遊戲房
	CreateGameRoom(ctx context.Context, hostAccountID int, maxPlayers int) (*entity.Game, error)

	// JoinGameRoom 加入遊戲房
	JoinGameRoom(ctx context.Context, gameID int, accountID int) error

	// StartGame 開始遊戲
	StartGame(ctx context.Context, gameID int, accountID int) (*entity.Game, error)

	// GetGameRoom 取得遊戲房資訊
	GetGameRoom(ctx context.Context, gameID int) (*entity.Game, error)
}

// gameRoomUseCase 遊戲房用例實作
type gameRoomUseCase struct {
	gameRepo       repository.GameRepository
	gamePlayerRepo repository.GamePlayerRepository
	accountRepo    repository.AccountRepository
	playerUseCase  PlayerUseCase
	gameUseCase    GameUseCase
}

// GameRoomUseCaseOptions 遊戲房用例選項
type GameRoomUseCaseOptions struct {
	GameRepo       repository.GameRepository
	GamePlayerRepo repository.GamePlayerRepository
	AccountRepo    repository.AccountRepository
	PlayerUseCase  PlayerUseCase
	GameUseCase    GameUseCase
}

// NewGameRoomUseCase 建立遊戲房用例
func NewGameRoomUseCase(opts *GameRoomUseCaseOptions) GameRoomUseCase {
	return &gameRoomUseCase{
		gameRepo:       opts.GameRepo,
		gamePlayerRepo: opts.GamePlayerRepo,
		accountRepo:    opts.AccountRepo,
		playerUseCase:  opts.PlayerUseCase,
		gameUseCase:    opts.GameUseCase,
	}
}

// CreateGameRoom 建立遊戲房
func (uc *gameRoomUseCase) CreateGameRoom(ctx context.Context, hostAccountID int, maxPlayers int) (*entity.Game, error) {
	// 驗證帳號是否存在
	_, err := uc.accountRepo.GetAccountById(ctx, hostAccountID)
	if err != nil {
		return nil, err
	}

	// 建立遊戲房
	game := entity.NewGameRoom(hostAccountID, maxPlayers)
	game, err = uc.gameRepo.CreateGame(ctx, game)
	if err != nil {
		return nil, err
	}

	// 房主自動加入遊戲
	gamePlayer := entity.NewGamePlayer(game.ID, hostAccountID, 1)
	_, err = uc.gamePlayerRepo.CreateGamePlayer(ctx, gamePlayer)
	if err != nil {
		return nil, err
	}

	return game, nil
}

// JoinGameRoom 加入遊戲房
func (uc *gameRoomUseCase) JoinGameRoom(ctx context.Context, gameID int, accountID int) error {
	// 取得遊戲房
	game, err := uc.gameRepo.GetGameById(ctx, gameID)
	if err != nil {
		return entity.ErrGameNotFound
	}

	// 檢查是否可以加入
	if err := game.CanJoin(); err != nil {
		return err
	}

	// 檢查是否已經在遊戲中
	exists, err := uc.gamePlayerRepo.ExistsByGameIDAndAccountID(ctx, gameID, accountID)
	if err != nil {
		return err
	}
	if exists {
		return entity.ErrAlreadyInGame
	}

	// 取得目前玩家數量作為加入順序
	currentCount, err := uc.gamePlayerRepo.CountByGameID(ctx, gameID)
	if err != nil {
		return err
	}

	// 加入遊戲
	gamePlayer := entity.NewGamePlayer(gameID, accountID, currentCount+1)
	_, err = uc.gamePlayerRepo.CreateGamePlayer(ctx, gamePlayer)
	if err != nil {
		return err
	}

	// 更新遊戲房玩家數
	game.CurrentPlayers = currentCount + 1
	err = uc.gameRepo.UpdateGame(ctx, game)
	if err != nil {
		return err
	}

	return nil
}

// StartGame 開始遊戲
func (uc *gameRoomUseCase) StartGame(ctx context.Context, gameID int, accountID int) (*entity.Game, error) {
	// 取得遊戲房
	game, err := uc.gameRepo.GetGameById(ctx, gameID)
	if err != nil {
		return nil, entity.ErrGameNotFound
	}

	// 檢查是否可以開始
	if err := game.CanStart(accountID); err != nil {
		return nil, err
	}

	// 產生遊戲 Token
	token, err := generateGameToken(256)
	if err != nil {
		return nil, err
	}
	game.Token = token

	// 取得所有加入遊戲的玩家（GamePlayer）
	gamePlayers, err := uc.gamePlayerRepo.GetGamePlayersByGameID(ctx, gameID)
	if err != nil {
		return nil, err
	}

	// 初始化身份卡
	identityCards := uc.playerUseCase.InitIdentityCards(len(gamePlayers))
	if len(identityCards) == 0 {
		return nil, entity.ErrPlayerCountNotSupported
	}

	// 為每個 GamePlayer 建立 Player 並分配身份
	var firstPlayerID int
	for i, gp := range gamePlayers {
		account, err := uc.accountRepo.GetAccountById(ctx, gp.AccountID)
		if err != nil {
			return nil, err
		}

		player, err := uc.playerUseCase.CreatePlayer(ctx, entity.NewPlayer(
			account.Name,
			gameID,
			identityCards[i],
			gp.JoinOrder,
		))
		if err != nil {
			return nil, err
		}

		// 第一位加入的玩家為第一位行動玩家
		if gp.JoinOrder == 1 {
			firstPlayerID = player.ID
		}
	}

	// 初始化牌堆
	if err := uc.gameUseCase.InitDeck(ctx, game); err != nil {
		return nil, err
	}

	// 為所有玩家發牌（每人 3 張）
	if err := uc.gameUseCase.DrawCardsForAllPlayers(ctx, game); err != nil {
		return nil, err
	}

	// 設定第一位行動玩家並更新遊戲狀態
	game.CurrentPlayerID = firstPlayerID
	game.Status = entity.GameRoomStatusPlaying
	err = uc.gameRepo.UpdateGame(ctx, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

// generateGameToken 產生遊戲 Token
func generateGameToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := cryptorand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetGameRoom 取得遊戲房資訊
func (uc *gameRoomUseCase) GetGameRoom(ctx context.Context, gameID int) (*entity.Game, error) {
	return uc.gameRepo.GetGameById(ctx, gameID)
}
