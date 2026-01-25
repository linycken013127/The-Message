package usecase

import (
	"context"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
)

// PassResult 跳過行動結果
type PassResult struct {
	AllPassed bool
	NewPhase  string
}

// ActionPhaseUseCase 行動階段用例介面
type ActionPhaseUseCase interface {
	// DrawCards 當前玩家抽牌（抽 2 張）
	DrawCards(ctx context.Context, gameID int, playerID int) ([]*entity.Card, error)
	// PlayFunctionCard 出功能牌
	PlayFunctionCard(ctx context.Context, gameID int, playerID int, cardID int) (*entity.Card, error)
	// Pass 跳過行動
	Pass(ctx context.Context, gameID int, playerID int) (*PassResult, error)
	// GetCurrentRound 取得當前回合
	GetCurrentRound(ctx context.Context, gameID int) (int, error)
}

// actionPhaseUseCase 行動階段用例實作
type actionPhaseUseCase struct {
	gameRepo       repository.GameRepository
	playerRepo     repository.PlayerRepository
	playerCardRepo repository.PlayerCardRepository
	deckUseCase    DeckUseCase
	actionPassRepo repository.ActionPassRepository
}

// ActionPhaseUseCaseOptions 行動階段用例選項
type ActionPhaseUseCaseOptions struct {
	GameRepo       repository.GameRepository
	PlayerRepo     repository.PlayerRepository
	PlayerCardRepo repository.PlayerCardRepository
	DeckUseCase    DeckUseCase
	ActionPassRepo repository.ActionPassRepository
}

// NewActionPhaseUseCase 建立行動階段用例
func NewActionPhaseUseCase(opts *ActionPhaseUseCaseOptions) ActionPhaseUseCase {
	return &actionPhaseUseCase{
		gameRepo:       opts.GameRepo,
		playerRepo:     opts.PlayerRepo,
		playerCardRepo: opts.PlayerCardRepo,
		deckUseCase:    opts.DeckUseCase,
		actionPassRepo: opts.ActionPassRepo,
	}
}

// DrawCards 當前玩家抽牌（抽 2 張）
func (uc *actionPhaseUseCase) DrawCards(ctx context.Context, gameID int, playerID int) ([]*entity.Card, error) {
	// 取得遊戲
	game, err := uc.gameRepo.GetGameById(ctx, gameID)
	if err != nil {
		return nil, entity.ErrGameNotFound
	}

	// 檢查遊戲是否在進行中
	if game.Status != entity.GameRoomStatusPlaying {
		return nil, entity.ErrGameNotFound
	}

	// 檢查是否在行動階段
	if game.Phase != entity.GamePhaseAction {
		return nil, entity.ErrNotInActionPhase
	}

	// 檢查是否是當前玩家
	player, err := uc.playerRepo.GetPlayerById(ctx, playerID)
	if err != nil {
		return nil, err
	}
	if game.CurrentPlayerID != playerID {
		return nil, entity.ErrNotYourTurn
	}

	// 取得牌堆
	decks, err := uc.deckUseCase.GetDecksByGameId(ctx, gameID)
	if err != nil {
		return nil, err
	}

	// 決定抽牌數量（正常抽 2 張，牌堆不足時抽剩餘的）
	drawCount := 2
	if len(decks) < 2 {
		drawCount = len(decks)
	}

	if drawCount == 0 {
		return nil, entity.ErrDeckNotEnough
	}

	// 抽牌
	drawnCards := make([]*entity.Card, 0, drawCount)
	for i := 0; i < drawCount; i++ {
		deck := decks[i]
		// 建立玩家手牌
		playerCard := entity.NewPlayerCard(playerID, gameID, deck.CardID, entity.PlayerCardTypeHand)
		_, err := uc.playerCardRepo.CreatePlayerCard(ctx, playerCard)
		if err != nil {
			return nil, err
		}
		// 從牌堆移除
		err = uc.deckUseCase.DeleteDeckFromGame(ctx, deck.ID)
		if err != nil {
			return nil, err
		}
		drawnCards = append(drawnCards, &entity.Card{ID: deck.CardID})
	}

	// 清除該玩家在本回合的 pass 記錄（如果有的話）
	round, _ := uc.GetCurrentRound(ctx, gameID)
	_ = uc.actionPassRepo.DeleteActionPassesByGameIDAndRound(ctx, gameID, round)

	_ = player // 避免 unused variable

	return drawnCards, nil
}

// PlayFunctionCard 出功能牌
func (uc *actionPhaseUseCase) PlayFunctionCard(ctx context.Context, gameID int, playerID int, cardID int) (*entity.Card, error) {
	// 取得遊戲
	game, err := uc.gameRepo.GetGameById(ctx, gameID)
	if err != nil {
		return nil, entity.ErrGameNotFound
	}

	// 檢查遊戲是否在進行中
	if game.Status != entity.GameRoomStatusPlaying {
		return nil, entity.ErrGameNotFound
	}

	// 檢查是否在行動階段
	if game.Phase != entity.GamePhaseAction {
		return nil, entity.ErrNotInActionPhase
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

	// 計算手牌數量，不能打出最後一張
	handCount, err := uc.playerCardRepo.CountHandCardsByPlayerID(ctx, playerID)
	if err != nil {
		return nil, err
	}
	if handCount <= 1 {
		return nil, entity.ErrCannotPlayLastCard
	}

	// 刪除手牌
	_, err = uc.playerCardRepo.DeletePlayerCardByPlayerIdAndCardId(ctx, playerID, gameID, cardID)
	if err != nil {
		return nil, err
	}

	// 清除所有玩家的 pass 記錄（有人出牌後重置 pass）
	round, _ := uc.GetCurrentRound(ctx, gameID)
	_ = uc.actionPassRepo.DeleteActionPassesByGameIDAndRound(ctx, gameID, round)

	return &entity.Card{ID: cardID}, nil
}

// Pass 跳過行動
func (uc *actionPhaseUseCase) Pass(ctx context.Context, gameID int, playerID int) (*PassResult, error) {
	// 取得遊戲（含玩家）
	game, err := uc.gameRepo.GetGameWithPlayers(ctx, gameID)
	if err != nil {
		return nil, entity.ErrGameNotFound
	}

	// 檢查遊戲是否在進行中
	if game.Status != entity.GameRoomStatusPlaying {
		return nil, entity.ErrGameNotFound
	}

	// 檢查是否在行動階段
	if game.Phase != entity.GamePhaseAction {
		return nil, entity.ErrNotInActionPhase
	}

	// 檢查是否是當前玩家
	if game.CurrentPlayerID != playerID {
		return nil, entity.ErrNotYourTurn
	}

	// 取得當前回合
	round, _ := uc.GetCurrentRound(ctx, gameID)

	// 記錄 pass
	actionPass := entity.NewActionPass(gameID, playerID, round)
	_, err = uc.actionPassRepo.CreateActionPass(ctx, actionPass)
	if err != nil {
		return nil, err
	}

	// 檢查是否所有存活玩家都已 pass
	passCount, err := uc.actionPassRepo.CountActionPassesByGameIDAndRound(ctx, gameID, round)
	if err != nil {
		return nil, err
	}

	// 計算存活玩家數
	alivePlayers := 0
	for _, p := range game.Players {
		if p.Status != entity.PlayerStatusDead {
			alivePlayers++
		}
	}

	allPassed := passCount >= alivePlayers
	if allPassed {
		// 所有人都 pass，進入情報階段
		game.Phase = entity.GamePhaseIntelligence
		err = uc.gameRepo.UpdateGame(ctx, game)
		if err != nil {
			return nil, err
		}
		// 清除 pass 記錄
		_ = uc.actionPassRepo.DeleteActionPassesByGameIDAndRound(ctx, gameID, round)
		return &PassResult{AllPassed: true, NewPhase: entity.GamePhaseIntelligence}, nil
	}

	// 切換到下一位玩家
	nextPlayerID := uc.getNextPlayerID(game, playerID)
	game.CurrentPlayerID = nextPlayerID
	err = uc.gameRepo.UpdateGame(ctx, game)
	if err != nil {
		return nil, err
	}

	return &PassResult{AllPassed: false, NewPhase: entity.GamePhaseAction}, nil
}

// GetCurrentRound 取得當前回合（簡化實作，總是回傳 1）
func (uc *actionPhaseUseCase) GetCurrentRound(ctx context.Context, gameID int) (int, error) {
	// 簡化實作：目前不追蹤回合數，固定為 1
	return 1, nil
}

// getNextPlayerID 取得下一位玩家 ID
func (uc *actionPhaseUseCase) getNextPlayerID(game *entity.Game, currentPlayerID int) int {
	players := game.Players
	currentIndex := -1

	for i, p := range players {
		if p.ID == currentPlayerID {
			currentIndex = i
			break
		}
	}

	// 找下一位存活的玩家
	for i := 1; i <= len(players); i++ {
		nextIndex := (currentIndex + i) % len(players)
		if players[nextIndex].Status != entity.PlayerStatusDead {
			return players[nextIndex].ID
		}
	}

	// 找不到下一位（不應該發生）
	return currentPlayerID
}
