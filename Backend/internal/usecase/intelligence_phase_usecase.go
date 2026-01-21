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
	// AcceptIntelligence 接收情報
	AcceptIntelligence(ctx context.Context, gameID int, playerID int) (*AcceptIntelligenceResult, error)
	// RejectIntelligence 拒絕情報
	RejectIntelligence(ctx context.Context, gameID int, playerID int) (*RejectIntelligenceResult, error)
}

// AcceptIntelligenceResult 接收情報結果
type AcceptIntelligenceResult struct {
	PlayerID int
	CardID   int
}

// RejectIntelligenceResult 拒絕情報結果
type RejectIntelligenceResult struct {
	Transfer       *entity.IntelligenceTransfer
	AutoAccepted   bool   // 是否自動接收（回到發送者或直達被拒絕）
	AutoAcceptedBy int    // 自動接收的玩家 ID
	Message        string // 回傳訊息
}

// intelligencePhaseUseCase 情報階段用例實作
type intelligencePhaseUseCase struct {
	gameRepo                 repository.GameRepository
	playerRepo               repository.PlayerRepository
	playerCardRepo           repository.PlayerCardRepository
	cardRepo                 repository.CardRepository
	intelligenceTransferRepo repository.IntelligenceTransferRepository
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

// AcceptIntelligence 接收情報
func (uc *intelligencePhaseUseCase) AcceptIntelligence(ctx context.Context, gameID int, playerID int) (*AcceptIntelligenceResult, error) {
	// 取得正在傳遞的情報
	transfer, err := uc.intelligenceTransferRepo.GetActiveTransferByGameID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, entity.ErrNoIntelligenceInTransit
	}

	// 檢查是否是目標玩家
	if transfer.CurrentTargetPlayerID != playerID {
		return nil, entity.ErrNotIntelligenceTarget
	}

	// 完成情報傳遞
	return uc.completeIntelligenceTransfer(ctx, transfer, playerID)
}

// RejectIntelligence 拒絕情報
func (uc *intelligencePhaseUseCase) RejectIntelligence(ctx context.Context, gameID int, playerID int) (*RejectIntelligenceResult, error) {
	// 取得遊戲（含玩家）
	game, err := uc.gameRepo.GetGameWithPlayers(ctx, gameID)
	if err != nil {
		return nil, entity.ErrGameNotFound
	}

	// 取得正在傳遞的情報
	transfer, err := uc.intelligenceTransferRepo.GetActiveTransferByGameID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, entity.ErrNoIntelligenceInTransit
	}

	// 檢查是否是目標玩家
	if transfer.CurrentTargetPlayerID != playerID {
		return nil, entity.ErrNotIntelligenceTarget
	}

	// 取得卡片資訊以判斷情報類型
	card, err := uc.cardRepo.GetCardById(ctx, transfer.CardID)
	if err != nil {
		return nil, err
	}

	// 判斷情報類型處理邏輯
	switch card.IntelligenceType {
	case entity.IntelligenceTypeDirect:
		// 直達情報被拒絕，自動歸屬發送者
		_, err := uc.completeIntelligenceTransfer(ctx, transfer, transfer.SenderPlayerID)
		if err != nil {
			return nil, err
		}
		return &RejectIntelligenceResult{
			Transfer:       transfer,
			AutoAccepted:   true,
			AutoAcceptedBy: transfer.SenderPlayerID,
			Message:        "情報已被發送者自動接收",
		}, nil

	case entity.IntelligenceTypeSecretTelegram, entity.IntelligenceTypeDocument:
		// 密電或文件，傳給下一位玩家
		nextTargetID := uc.getNextAlivePlayerID(game, playerID)

		// 檢查是否回到發送者
		if nextTargetID == transfer.SenderPlayerID {
			// 回到發送者，自動接收
			_, err := uc.completeIntelligenceTransfer(ctx, transfer, transfer.SenderPlayerID)
			if err != nil {
				return nil, err
			}
			return &RejectIntelligenceResult{
				Transfer:       transfer,
				AutoAccepted:   true,
				AutoAcceptedBy: transfer.SenderPlayerID,
				Message:        "情報已被發送者自動接收",
			}, nil
		}

		// 更新目標玩家
		transfer.UpdateTarget(nextTargetID)
		err = uc.intelligenceTransferRepo.UpdateIntelligenceTransfer(ctx, transfer)
		if err != nil {
			return nil, err
		}

		// 更新當前玩家為新目標
		game.CurrentPlayerID = nextTargetID
		err = uc.gameRepo.UpdateGame(ctx, game)
		if err != nil {
			return nil, err
		}

		// 附加卡片資訊
		transfer.Card = card

		return &RejectIntelligenceResult{
			Transfer:     transfer,
			AutoAccepted: false,
			Message:      "情報已拒絕",
		}, nil

	default:
		// 未知類型，預設為密電處理
		nextTargetID := uc.getNextAlivePlayerID(game, playerID)
		if nextTargetID == transfer.SenderPlayerID {
			_, err := uc.completeIntelligenceTransfer(ctx, transfer, transfer.SenderPlayerID)
			if err != nil {
				return nil, err
			}
			return &RejectIntelligenceResult{
				Transfer:       transfer,
				AutoAccepted:   true,
				AutoAcceptedBy: transfer.SenderPlayerID,
				Message:        "情報已被發送者自動接收",
			}, nil
		}

		transfer.UpdateTarget(nextTargetID)
		err = uc.intelligenceTransferRepo.UpdateIntelligenceTransfer(ctx, transfer)
		if err != nil {
			return nil, err
		}

		game.CurrentPlayerID = nextTargetID
		err = uc.gameRepo.UpdateGame(ctx, game)
		if err != nil {
			return nil, err
		}

		transfer.Card = card
		return &RejectIntelligenceResult{
			Transfer:     transfer,
			AutoAccepted: false,
			Message:      "情報已拒絕",
		}, nil
	}
}

// completeIntelligenceTransfer 完成情報傳遞
func (uc *intelligencePhaseUseCase) completeIntelligenceTransfer(ctx context.Context, transfer *entity.IntelligenceTransfer, receiverID int) (*AcceptIntelligenceResult, error) {
	// 取得遊戲
	game, err := uc.gameRepo.GetGameWithPlayers(ctx, transfer.GameID)
	if err != nil {
		return nil, err
	}

	// 將卡片加入接收者的情報區
	intelligenceCard := entity.NewPlayerCard(receiverID, transfer.GameID, transfer.CardID, entity.PlayerCardTypeIntelligence)
	_, err = uc.playerCardRepo.CreatePlayerCard(ctx, intelligenceCard)
	if err != nil {
		return nil, err
	}

	// 標記情報傳遞為已完成
	transfer.Complete()
	err = uc.intelligenceTransferRepo.UpdateIntelligenceTransfer(ctx, transfer)
	if err != nil {
		return nil, err
	}

	// 刪除情報傳遞記錄（或保留作為歷史記錄）
	err = uc.intelligenceTransferRepo.DeleteIntelligenceTransfer(ctx, transfer.ID)
	if err != nil {
		return nil, err
	}

	// 遊戲進入下一回合（行動階段）
	game.Phase = entity.GamePhaseAction
	game.CurrentPlayerID = receiverID
	err = uc.gameRepo.UpdateGame(ctx, game)
	if err != nil {
		return nil, err
	}

	return &AcceptIntelligenceResult{
		PlayerID: receiverID,
		CardID:   transfer.CardID,
	}, nil
}
