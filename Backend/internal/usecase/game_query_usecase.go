package usecase

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
)

// GameQueryUseCase 遊戲查詢用例介面
type GameQueryUseCase interface {
	// GetGameInfo 取得遊戲資訊
	GetGameInfo(ctx context.Context, gameID int, accountID int) (*GameInfoResult, error)
}

// GameInfoResult 遊戲資訊結果
type GameInfoResult struct {
	GameID            int                `json:"gameId"`
	Status            string             `json:"status"`
	Phase             string             `json:"phase"`
	CurrentPlayerID   int                `json:"currentPlayerId"`
	HasIntelInTransit bool               `json:"hasIntelInTransit"`
	MyFaction         string             `json:"myFaction"`
	MyRedCount        int                `json:"myRedCount"`
	MyBlueCount       int                `json:"myBlueCount"`
	MyBlackCount      int                `json:"myBlackCount"`
	MyHandCardCount   int                `json:"myHandCardCount"`
	Players           []PlayerPublicInfo `json:"players"`
}

// PlayerPublicInfo 玩家公開資訊
type PlayerPublicInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Alive      bool   `json:"alive"`
	IntelCount int    `json:"intelCount"`
}

// gameQueryUseCase 遊戲查詢用例實作
type gameQueryUseCase struct {
	gameRepo                 repository.GameRepository
	playerRepo               repository.PlayerRepository
	playerCardRepo           repository.PlayerCardRepository
	gamePlayerRepo           repository.GamePlayerRepository
	intelligenceTransferRepo repository.IntelligenceTransferRepository
}

// GameQueryUseCaseOptions 遊戲查詢用例選項
type GameQueryUseCaseOptions struct {
	GameRepo                 repository.GameRepository
	PlayerRepo               repository.PlayerRepository
	PlayerCardRepo           repository.PlayerCardRepository
	GamePlayerRepo           repository.GamePlayerRepository
	IntelligenceTransferRepo repository.IntelligenceTransferRepository
}

// NewGameQueryUseCase 建立遊戲查詢用例
func NewGameQueryUseCase(opts *GameQueryUseCaseOptions) GameQueryUseCase {
	return &gameQueryUseCase{
		gameRepo:                 opts.GameRepo,
		playerRepo:               opts.PlayerRepo,
		playerCardRepo:           opts.PlayerCardRepo,
		gamePlayerRepo:           opts.GamePlayerRepo,
		intelligenceTransferRepo: opts.IntelligenceTransferRepo,
	}
}

// GetGameInfo 取得遊戲資訊
func (uc *gameQueryUseCase) GetGameInfo(ctx context.Context, gameID int, accountID int) (*GameInfoResult, error) {
	// 取得遊戲
	game, err := uc.gameRepo.GetGameWithPlayers(ctx, gameID)
	if err != nil {
		return nil, entity.ErrGameNotFound
	}

	// 檢查是否有情報正在傳遞
	hasIntelInTransit := false
	activeTransfer, err := uc.intelligenceTransferRepo.GetActiveTransferByGameID(ctx, gameID)
	if err == nil && activeTransfer != nil {
		hasIntelInTransit = true
	}

	// 建立結果
	result := &GameInfoResult{
		GameID:            game.ID,
		Status:            game.Status,
		Phase:             game.Phase,
		CurrentPlayerID:   game.CurrentPlayerID,
		HasIntelInTransit: hasIntelInTransit,
		Players:           make([]PlayerPublicInfo, 0),
	}

	// 取得所有玩家並計算公開資訊
	players, err := uc.playerRepo.GetPlayersByGameId(ctx, gameID)
	if err != nil {
		return nil, err
	}

	// 找出請求者對應的玩家
	gamePlayer, err := uc.gamePlayerRepo.GetGamePlayerByGameIDAndAccountID(ctx, gameID, accountID)
	var myPlayerID int
	if err == nil && gamePlayer != nil {
		// 根據加入順序找出對應的玩家
		for _, p := range players {
			if p.OrderNumber == gamePlayer.JoinOrder {
				myPlayerID = p.ID
				break
			}
		}
	}

	// 處理每個玩家
	for _, player := range players {
		// 計算情報數量
		playerWithCards, err := uc.playerRepo.GetPlayerWithPlayerCards(ctx, player.ID)
		if err != nil {
			continue
		}

		intelCount := 0
		var redCount, blueCount, blackCount, handCardCount int
		for _, pc := range playerWithCards.PlayerCards {
			if pc.Type == entity.PlayerCardTypeIntelligence {
				intelCount++
				switch pc.Card.Color {
				case entity.CardColorRed:
					redCount++
				case entity.CardColorBlue:
					blueCount++
				case entity.CardColorBlack:
					blackCount++
				}
			} else if pc.Type == entity.PlayerCardTypeHand {
				handCardCount++
			}
		}

		// 新增玩家公開資訊
		result.Players = append(result.Players, PlayerPublicInfo{
			ID:         player.ID,
			Name:       player.Name,
			Alive:      player.Status == entity.PlayerStatusAlive,
			IntelCount: intelCount,
		})

		// 如果是請求者，新增個人資訊
		if player.ID == myPlayerID {
			result.MyFaction = player.IdentityCard
			result.MyRedCount = redCount
			result.MyBlueCount = blueCount
			result.MyBlackCount = blackCount
			result.MyHandCardCount = handCardCount
		}
	}

	return result, nil
}
