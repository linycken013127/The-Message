package usecase

import (
	"context"
	"errors"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
)

// PlayerInfo 玩家資訊（用於建立遊戲）
type PlayerInfo struct {
	ID   string
	Name string
}

// CreateGameRequest 建立遊戲請求
type CreateGameRequest struct {
	Players []PlayerInfo
}

// AcceptCardResult 接收卡片結果
type AcceptCardResult struct {
	Accepted bool
	Winner   *entity.Player
}

// PlayerUseCase 玩家用例介面
type PlayerUseCase interface {
	// InitPlayers 初始化玩家
	InitPlayers(ctx context.Context, game *entity.Game, req CreateGameRequest) error
	// InitIdentityCards 初始化身份卡
	InitIdentityCards(playersCount int) []string
	// CanPlayCard 檢查玩家是否可以出牌
	CanPlayCard(ctx context.Context, player *entity.Player) (bool, error)
	// CheckPlayerCardExist 檢查玩家手牌是否存在
	CheckPlayerCardExist(ctx context.Context, playerID int, gameID int, cardID int) (bool, error)
	// CreatePlayer 建立玩家
	CreatePlayer(ctx context.Context, player *entity.Player) (*entity.Player, error)
	// CreatePlayerCard 建立玩家手牌
	CreatePlayerCard(ctx context.Context, card *entity.PlayerCard) error
	// GetPlayerById 根據 ID 取得玩家
	GetPlayerById(ctx context.Context, id int) (*entity.Player, error)
	// GetPlayersByGameId 根據遊戲 ID 取得玩家
	GetPlayersByGameId(ctx context.Context, gameID int) ([]*entity.Player, error)
	// GetHandCardId 取得玩家手牌
	GetHandCardId(player *entity.Player, cardID int) (*entity.PlayerCard, error)
	// PlayCard 出牌
	PlayCard(ctx context.Context, playerID int, cardID int) (*entity.Game, *entity.Card, error)
	// TransmitIntelligence 傳遞情報（包含驗證）
	TransmitIntelligence(ctx context.Context, playerID int, cardID int) (bool, error)
	// TransmitIntelligenceCard 傳情報卡
	TransmitIntelligenceCard(ctx context.Context, playerID int, gameID int, cardID int) (bool, error)
	// AcceptCard 接收卡片
	AcceptCard(ctx context.Context, playerID int, accept bool) (bool, error)
	// AcceptCardAndCheckWin 接收卡片並檢查勝利
	AcceptCardAndCheckWin(ctx context.Context, playerID int, accept bool) (*AcceptCardResult, error)
	// CheckWin 檢查勝利
	CheckWin(ctx context.Context, playerID int) (*entity.Player, error)
	// GetPlayerWithPlayerCards 取得玩家及其手牌
	GetPlayerWithPlayerCards(ctx context.Context, playerID int) (*entity.Player, error)
	// SetGameUseCase 設定遊戲用例（處理循環依賴）
	SetGameUseCase(gameUseCase GameUseCase)
}

// playerUseCase 玩家用例實作
type playerUseCase struct {
	playerRepo       repository.PlayerRepository
	playerCardRepo   repository.PlayerCardRepository
	gameRepo         repository.GameRepository
	gameProgressRepo repository.GameProgressRepository
	gameUseCase      GameUseCase
}

// PlayerUseCaseOptions 玩家用例選項
type PlayerUseCaseOptions struct {
	PlayerRepo       repository.PlayerRepository
	PlayerCardRepo   repository.PlayerCardRepository
	GameRepo         repository.GameRepository
	GameProgressRepo repository.GameProgressRepository
}

// NewPlayerUseCase 建立玩家用例
func NewPlayerUseCase(opts *PlayerUseCaseOptions) PlayerUseCase {
	return &playerUseCase{
		playerRepo:       opts.PlayerRepo,
		playerCardRepo:   opts.PlayerCardRepo,
		gameRepo:         opts.GameRepo,
		gameProgressRepo: opts.GameProgressRepo,
	}
}

// SetGameUseCase 設定遊戲用例
func (uc *playerUseCase) SetGameUseCase(gameUseCase GameUseCase) {
	uc.gameUseCase = gameUseCase
}

// InitPlayers 初始化玩家
func (uc *playerUseCase) InitPlayers(ctx context.Context, game *entity.Game, req CreateGameRequest) error {
	identityCards := uc.InitIdentityCards(len(req.Players))
	for i, reqPlayer := range req.Players {
		_, err := uc.CreatePlayer(ctx, entity.NewPlayer(
			reqPlayer.Name,
			game.ID,
			identityCards[i],
			i+1,
		))
		if err != nil {
			return err
		}
	}
	return nil
}

// InitIdentityCards 初始化身份卡
// 根據玩家數量分配身份：
// 3人: 潛伏1, 軍情1, 打醬油1
// 5人: 潛伏2, 軍情2, 打醬油1
// 6人: 潛伏2, 軍情2, 打醬油2
// 7人: 潛伏2, 軍情2, 打醬油3
// 8人: 潛伏3, 軍情3, 打醬油2
// 9人: 潛伏3, 軍情3, 打醬油3
func (uc *playerUseCase) InitIdentityCards(playersCount int) []string {
	identityCards := make([]string, 0, playersCount)

	var undercover, military, bystander int

	switch playersCount {
	case 3:
		undercover, military, bystander = 1, 1, 1
	case 5:
		undercover, military, bystander = 2, 2, 1
	case 6:
		undercover, military, bystander = 2, 2, 2
	case 7:
		undercover, military, bystander = 2, 2, 3
	case 8:
		undercover, military, bystander = 3, 3, 2
	case 9:
		undercover, military, bystander = 3, 3, 3
	default:
		// 不支援的人數，返回空陣列
		return identityCards
	}

	// 加入潛伏戰線
	for i := 0; i < undercover; i++ {
		identityCards = append(identityCards, entity.IdentityUndercoverFront)
	}
	// 加入軍情處
	for i := 0; i < military; i++ {
		identityCards = append(identityCards, entity.IdentityMilitaryAgency)
	}
	// 加入打醬油
	for i := 0; i < bystander; i++ {
		identityCards = append(identityCards, entity.IdentityBystander)
	}

	identityCards = ShuffleIdentityCards(identityCards)
	return identityCards
}

// CanPlayCard 檢查玩家是否可以出牌
func (uc *playerUseCase) CanPlayCard(ctx context.Context, player *entity.Player) (bool, error) {
	if player.Game.Status == entity.GameStatusEnd {
		return false, errors.New("遊戲已結束")
	}

	if player.Status == entity.PlayerStatusDead {
		return false, errors.New("你已死亡")
	}

	if player.Game.CurrentPlayerID != player.ID {
		return false, errors.New("尚未輪到你出牌")
	}

	return true, nil
}

// CheckPlayerCardExist 檢查玩家手牌是否存在
func (uc *playerUseCase) CheckPlayerCardExist(ctx context.Context, playerID int, gameID int, cardID int) (bool, error) {
	exist, err := uc.playerCardRepo.ExistPlayerCardByPlayerIdAndCardId(ctx, playerID, gameID, cardID)
	if err != nil {
		return false, err
	}
	return exist, nil
}

// CreatePlayer 建立玩家
func (uc *playerUseCase) CreatePlayer(ctx context.Context, player *entity.Player) (*entity.Player, error) {
	player, err := uc.playerRepo.CreatePlayer(ctx, player)
	if err != nil {
		return nil, err
	}
	return player, nil
}

// CreatePlayerCard 建立玩家手牌
func (uc *playerUseCase) CreatePlayerCard(ctx context.Context, card *entity.PlayerCard) error {
	_, err := uc.playerCardRepo.CreatePlayerCard(ctx, card)
	return err
}

// GetPlayerById 根據 ID 取得玩家
func (uc *playerUseCase) GetPlayerById(ctx context.Context, id int) (*entity.Player, error) {
	player, err := uc.playerRepo.GetPlayerById(ctx, id)
	if err != nil {
		return nil, err
	}
	return player, nil
}

// GetPlayersByGameId 根據遊戲 ID 取得玩家
func (uc *playerUseCase) GetPlayersByGameId(ctx context.Context, gameID int) ([]*entity.Player, error) {
	players, err := uc.playerRepo.GetPlayersByGameId(ctx, gameID)
	if err != nil {
		return nil, err
	}
	return players, nil
}

// GetPlayerWithPlayerCards 取得玩家及其手牌
func (uc *playerUseCase) GetPlayerWithPlayerCards(ctx context.Context, playerID int) (*entity.Player, error) {
	return uc.playerRepo.GetPlayerWithPlayerCards(ctx, playerID)
}

// GetHandCardId 取得玩家手牌
func (uc *playerUseCase) GetHandCardId(player *entity.Player, cardID int) (*entity.PlayerCard, error) {
	for _, card := range player.PlayerCards {
		if card.CardID == cardID && card.Type == entity.PlayerCardTypeHand {
			return &card, nil
		}
	}
	return nil, errors.New("找不到手牌")
}

// PlayCard 出牌
func (uc *playerUseCase) PlayCard(ctx context.Context, playerID int, cardID int) (*entity.Game, *entity.Card, error) {
	player, err := uc.playerRepo.GetPlayerWithGamePlayersAndPlayerCardsCard(ctx, playerID)
	if err != nil {
		return nil, nil, err
	}

	result, err := uc.CanPlayCard(ctx, player)
	if !result || err != nil {
		return nil, nil, err
	}

	handCard, err := uc.GetHandCardId(player, cardID)
	if err != nil {
		return nil, nil, err
	}

	game, err := uc.gameUseCase.NextPlayer(ctx, player)
	if err != nil {
		return nil, nil, err
	}

	err = uc.playerCardRepo.DeletePlayerCard(ctx, handCard.ID)
	if err != nil {
		return nil, nil, err
	}

	err = uc.gameRepo.UpdateGame(ctx, game)
	if err != nil {
		return nil, nil, err
	}

	return game, &handCard.Card, nil
}

// TransmitIntelligence 傳遞情報（包含驗證）
func (uc *playerUseCase) TransmitIntelligence(ctx context.Context, playerID int, cardID int) (bool, error) {
	player, err := uc.GetPlayerById(ctx, playerID)
	if err != nil || player == nil {
		return false, errors.New("Player not found")
	}

	exist, err := uc.CheckPlayerCardExist(ctx, playerID, player.GameID, cardID)
	if err != nil {
		return false, err
	}
	if !exist {
		return false, errors.New("Card not found")
	}

	return uc.TransmitIntelligenceCard(ctx, playerID, player.GameID, cardID)
}

// TransmitIntelligenceCard 傳情報卡
func (uc *playerUseCase) TransmitIntelligenceCard(ctx context.Context, playerID int, gameID int, cardID int) (bool, error) {
	player, err := uc.playerRepo.GetPlayerWithGamePlayersAndPlayerCardsCard(ctx, playerID)
	if err != nil {
		return false, err
	}

	result, err := uc.CanPlayCard(ctx, player)
	if !result || err != nil {
		return false, err
	}

	game, err := uc.gameUseCase.NextPlayer(ctx, player)
	if err != nil {
		return false, err
	}

	ret, err := uc.playerCardRepo.DeletePlayerCardByPlayerIdAndCardId(ctx, playerID, gameID, cardID)
	if err != nil {
		return false, err
	}

	err = uc.gameRepo.UpdateGame(ctx, game)
	if err != nil {
		return false, err
	}

	_, err = uc.gameProgressRepo.CreateGameProgress(ctx, entity.NewGameProgress(
		playerID,
		game.ID,
		cardID,
		entity.ActionTransmitIntelligence,
		game.CurrentPlayerID,
	))
	if err != nil {
		return false, err
	}

	return ret, nil
}

// AcceptCard 接收卡片
func (uc *playerUseCase) AcceptCard(ctx context.Context, playerID int, accept bool) (bool, error) {
	player, err := uc.playerRepo.GetPlayerWithGamePlayersAndPlayerCardsCard(ctx, playerID)
	if err != nil {
		return false, err
	}

	result, err := uc.CanPlayCard(ctx, player)
	if !result || err != nil {
		return false, err
	}

	game, err := uc.gameUseCase.NextPlayer(ctx, player)
	if err != nil {
		return false, err
	}

	gameID := game.ID
	gameProgress, err := uc.gameProgressRepo.GetGameProgress(ctx, playerID, gameID)
	if err != nil {
		return false, err
	}
	cardID := gameProgress.CardID

	res := accept
	if accept {
		_, err := uc.playerCardRepo.CreatePlayerCard(ctx, entity.NewPlayerCard(
			playerID,
			gameID,
			cardID,
			entity.PlayerCardTypeIntelligence,
		))
		if err != nil {
			return false, err
		}
		uc.gameUseCase.UpdateStatus(ctx, game, entity.GameStatusActionCardStage)
	} else {
		_, err := uc.gameProgressRepo.UpdateGameProgress(ctx, gameProgress, game.CurrentPlayerID)
		if err != nil {
			return false, err
		}

		err = uc.gameRepo.UpdateGame(ctx, game)
		if err != nil {
			return false, err
		}
	}

	return res, nil
}

// AcceptCardAndCheckWin 接收卡片並檢查勝利
func (uc *playerUseCase) AcceptCardAndCheckWin(ctx context.Context, playerID int, accept bool) (*AcceptCardResult, error) {
	accepted, err := uc.AcceptCard(ctx, playerID, accept)
	if err != nil {
		return nil, err
	}

	winner, _ := uc.CheckWin(ctx, playerID)

	return &AcceptCardResult{
		Accepted: accepted,
		Winner:   winner,
	}, nil
}

// CheckWin 檢查勝利
func (uc *playerUseCase) CheckWin(ctx context.Context, playerID int) (*entity.Player, error) {
	player, err := uc.playerRepo.GetPlayerWithGamePlayersAndPlayerCardsCard(ctx, playerID)
	if err != nil {
		return nil, err
	}

	win := 0
	var winPlayer *entity.Player
	for _, p := range player.Game.Players {
		win = 0
		for _, card := range p.PlayerCards {
			if card.Type == entity.PlayerCardTypeIntelligence && p.IdentityCard == entity.IdentityMilitaryAgency && card.Card.Color == entity.CardColorRed {
				win++
				if win == 3 {
					winPlayer = &p
					break
				}
			}

			if card.Type == entity.PlayerCardTypeIntelligence && p.IdentityCard == entity.IdentityUndercoverFront && card.Card.Color == entity.CardColorBlue {
				win++
				if win == 3 {
					winPlayer = &p
					break
				}
			}

			if card.Type == entity.PlayerCardTypeIntelligence && p.IdentityCard == entity.IdentityMilitaryAgency && (card.Card.Color == entity.CardColorRed || card.Card.Color == entity.CardColorBlue) {
				win++
				if win == 5 {
					winPlayer = &p
					break
				}
			}
		}
	}
	return winPlayer, nil
}
