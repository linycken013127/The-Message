package usecase

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
)

// IntelligencePhaseUseCase 情報階段用例介面
type IntelligencePhaseUseCase interface {
	// PassIntelligenceCard 傳遞情報牌
	PassIntelligenceCard(ctx context.Context, gameID int, playerID int, cardID int, targetPlayerID int) (*entity.IntelligenceTransfer, error)
	// GetActiveTransfer 取得正在傳遞的情報
	GetActiveTransfer(ctx context.Context, gameID int) (*entity.IntelligenceTransfer, error)
}

// intelligencePhaseUseCase 情報階段用例實作
type intelligencePhaseUseCase struct {
	gameRepo                   repository.GameRepository
	playerRepo                 repository.PlayerRepository
	playerCardRepo             repository.PlayerCardRepository
	cardRepo                   repository.CardRepository
	intelligenceTransferRepo   repository.IntelligenceTransferRepository
}

// IntelligencePhaseUseCaseOptions 情報階段用例選項
type IntelligencePhaseUseCaseOptions struct {
	GameRepo                 repository.GameRepository
	PlayerRepo               repository.PlayerRepository
	PlayerCardRepo           repository.PlayerCardRepository
	CardRepo                 repository.CardRepository
	IntelligenceTransferRepo repository.IntelligenceTransferRepository
}

// NewIntelligencePhaseUseCase 建立情報階段用例
func NewIntelligencePhaseUseCase(opts *IntelligencePhaseUseCaseOptions) IntelligencePhaseUseCase {
	return &intelligencePhaseUseCase{
		gameRepo:                 opts.GameRepo,
		playerRepo:               opts.PlayerRepo,
		playerCardRepo:           opts.PlayerCardRepo,
		cardRepo:                 opts.CardRepo,
		intelligenceTransferRepo: opts.IntelligenceTransferRepo,
	}
}

// PassIntelligenceCard 傳遞情報牌
func (uc *intelligencePhaseUseCase) PassIntelligenceCard(ctx context.Context, gameID int, playerID int, cardID int, targetPlayerID int) (*entity.IntelligenceTransfer, error) {
	// 取得遊戲（含玩家）
	game, err := uc.gameRepo.GetGameWithPlayers(ctx, gameID)
	if err != nil {
		return nil, entity.ErrGameNotFound
	}

	// 檢查遊戲是否在進行中
	if game.Status != entity.GameRoomStatusPlaying {
		return nil, entity.ErrGameNotFound
	}

	// 檢查是否在情報階段
	if game.Phase != entity.GamePhaseIntelligence {
		return nil, entity.ErrNotInIntelligencePhase
	}

	// 檢查是否是當前玩家
	if game.CurrentPlayerID != playerID {
		return nil, entity.ErrNotYourTurn
	}

	// 檢查手牌是否存在
	exists, err := uc.playerCardRepo.ExistPlayerCardByPlayerIdAndCardId(ctx, playerID, gameID, cardID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, entity.ErrCardNotInHand
	}

	// 取得卡片資訊
	card, err := uc.cardRepo.GetCardById(ctx, cardID)
	if err != nil {
		return nil, err
	}

	// 根據卡片類型決定傳遞方式
	var currentTargetID, originalTargetID int
	var faceUp bool

	switch card.IntelligenceType {
	case entity.IntelligenceTypeSecretTelegram:
		// 密電：蓋牌向右傳
		faceUp = false
		currentTargetID = uc.getNextAlivePlayerID(game, playerID)
		originalTargetID = currentTargetID

	case entity.IntelligenceTypeDocument:
		// 文件：明牌向右傳
		faceUp = true
		currentTargetID = uc.getNextAlivePlayerID(game, playerID)
		originalTargetID = currentTargetID

	case entity.IntelligenceTypeDirect:
		// 直達：蓋牌指定玩家
		if targetPlayerID == 0 {
			return nil, entity.ErrDirectCardNeedsTarget
		}
		if targetPlayerID == playerID {
			return nil, entity.ErrCannotTargetSelf
		}
		// 檢查目標玩家是否存活
		targetPlayer, err := uc.playerRepo.GetPlayerById(ctx, targetPlayerID)
		if err != nil {
			return nil, err
		}
		if targetPlayer.Status == entity.PlayerStatusDead {
			return nil, entity.ErrCannotTargetDeadPlayer
		}
		faceUp = false
		currentTargetID = targetPlayerID
		originalTargetID = targetPlayerID

	default:
		// 未知類型，預設為密電
		faceUp = false
		currentTargetID = uc.getNextAlivePlayerID(game, playerID)
		originalTargetID = currentTargetID
	}

	// 從手牌移除卡片
	_, err = uc.playerCardRepo.DeletePlayerCardByPlayerIdAndCardId(ctx, playerID, gameID, cardID)
	if err != nil {
		return nil, err
	}

	// 建立情報傳遞記錄
	transfer := entity.NewIntelligenceTransfer(
		gameID,
		cardID,
		playerID,
		currentTargetID,
		originalTargetID,
		faceUp,
	)

	transfer, err = uc.intelligenceTransferRepo.CreateIntelligenceTransfer(ctx, transfer)
	if err != nil {
		return nil, err
	}

	// 更新當前玩家為目標玩家
	game.CurrentPlayerID = currentTargetID
	err = uc.gameRepo.UpdateGame(ctx, game)
	if err != nil {
		return nil, err
	}

	// 附加卡片資訊
	transfer.Card = card

	return transfer, nil
}

// GetActiveTransfer 取得正在傳遞的情報
func (uc *intelligencePhaseUseCase) GetActiveTransfer(ctx context.Context, gameID int) (*entity.IntelligenceTransfer, error) {
	return uc.intelligenceTransferRepo.GetActiveTransferByGameID(ctx, gameID)
}

// getNextAlivePlayerID 取得下一位存活玩家 ID
func (uc *intelligencePhaseUseCase) getNextAlivePlayerID(game *entity.Game, currentPlayerID int) int {
	players := game.Players
	currentIndex := -1

	// 找到當前玩家的索引
	for i, p := range players {
		if p.ID == currentPlayerID {
			currentIndex = i
			break
		}
	}

	// 往右找下一位存活玩家
	for i := 1; i <= len(players); i++ {
		nextIndex := (currentIndex + i) % len(players)
		if players[nextIndex].Status != entity.PlayerStatusDead {
			return players[nextIndex].ID
		}
	}

	// 找不到（不應該發生）
	return currentPlayerID
}
