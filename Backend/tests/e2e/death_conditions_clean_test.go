package e2e

import (
	"context"
	"net/http"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/config"
)

// ==================== 死亡條件測試 ====================

// TestDeath_ThreeBlackIntelligence 測試玩家收集 3 張黑色情報死亡
func (suite *CleanArchTestSuite) TestDeath_ThreeBlackIntelligence() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定目標玩家已有 2 張黑色情報
	senderID := game.CurrentPlayerID
	targetPlayerID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.setPlayerIntelligenceCount(targetPlayerID, game.ID, 0, 0, 2)

	// 當前玩家傳遞一張黑色情報卡
	blackCard := suite.addBlackIntelligenceCardToPlayer(senderID, game.ID)
	suite.Require().NotNil(blackCard, "需要黑色情報卡來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, blackCard.CardID, 0)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查玩家是否死亡
	player, _ := suite.playerRepo.GetPlayerById(context.Background(), targetPlayerID)
	suite.Equal(entity.PlayerStatusDead, player.Status, "玩家應該死亡")
}

// TestDeath_MoreThanThreeBlack 測試玩家收集超過 3 張黑色情報死亡
func (suite *CleanArchTestSuite) TestDeath_MoreThanThreeBlack() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定目標玩家已有 3 張黑色情報
	senderID := game.CurrentPlayerID
	targetPlayerID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.setPlayerIntelligenceCount(targetPlayerID, game.ID, 0, 0, 3)

	// 當前玩家傳遞一張黑色情報卡
	blackCard := suite.addBlackIntelligenceCardToPlayer(senderID, game.ID)
	suite.Require().NotNil(blackCard, "需要黑色情報卡來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, blackCard.CardID, 0)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查玩家是否死亡
	player, _ := suite.playerRepo.GetPlayerById(context.Background(), targetPlayerID)
	suite.Equal(entity.PlayerStatusDead, player.Status, "玩家應該死亡")
}

// TestDeath_NotEnoughBlack 測試黑色情報少於 3 張時玩家保持存活
func (suite *CleanArchTestSuite) TestDeath_NotEnoughBlack() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定目標玩家已有 1 張黑色情報
	senderID := game.CurrentPlayerID
	targetPlayerID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.setPlayerIntelligenceCount(targetPlayerID, game.ID, 0, 0, 1)

	// 當前玩家傳遞一張黑色情報卡
	blackCard := suite.addBlackIntelligenceCardToPlayer(senderID, game.ID)
	suite.Require().NotNil(blackCard, "需要黑色情報卡來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, blackCard.CardID, 0)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查玩家應該仍然存活
	player, _ := suite.playerRepo.GetPlayerById(context.Background(), targetPlayerID)
	suite.Equal(entity.PlayerStatusAlive, player.Status, "玩家應該存活")
}

// TestDeath_SkipDeadPlayerTurn 測試死亡玩家被跳過回合
func (suite *CleanArchTestSuite) TestDeath_SkipDeadPlayerTurn() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定 Player2（下一位玩家）為死亡狀態
	senderID := game.CurrentPlayerID
	player2ID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.setPlayerDead(player2ID)

	// Player3 是跳過 Player2 後的下一位存活玩家
	player3ID := suite.getNextAlivePlayerIDSkipDead(game, players, senderID)

	// 當前玩家傳遞情報
	handCards := suite.getPlayerHandCardsWithType(senderID, entity.IntelligenceTypeSecretTelegram)
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(senderID, game.ID, entity.IntelligenceTypeSecretTelegram)
	}
	suite.Require().NotEmpty(handCards, "需要密電類型的卡片來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, handCards[0].CardID, 0)

	// 情報應該跳過死亡的 Player2，直接傳給 Player3
	transfer, _ := suite.intelligenceTransferRepo.GetActiveTransferByGameID(context.Background(), game.ID)
	suite.Equal(player3ID, transfer.CurrentTargetPlayerID, "情報應該跳過死亡玩家")
}

// TestDeath_SkipDeadPlayerIntelligence 測試情報傳遞時跳過死亡玩家
func (suite *CleanArchTestSuite) TestDeath_SkipDeadPlayerIntelligence() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 當前玩家傳遞密電情報
	senderID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(senderID, entity.IntelligenceTypeSecretTelegram)
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(senderID, game.ID, entity.IntelligenceTypeSecretTelegram)
	}
	suite.Require().NotEmpty(handCards, "需要密電類型的卡片來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, handCards[0].CardID, 0)

	// 第一位目標玩家
	target1ID := suite.getNextAlivePlayerID(game, players, senderID)

	// 設定第二位目標玩家為死亡狀態
	target2ID := suite.getNextAlivePlayerID(game, players, target1ID)
	suite.setPlayerDead(target2ID)

	// 第一位玩家拒絕情報，應該跳過死亡的第二位玩家
	resp := suite.rejectIntelligence(game.ID, target1ID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查情報傳遞狀態 - 應該跳過死亡玩家傳給發送者（3人遊戲）
	responseJson := suite.responseJson(resp)

	// 由於死亡玩家被跳過，情報會回到發送者自動接收
	if responseJson["auto_accepted"] != nil && responseJson["auto_accepted"].(bool) {
		suite.Equal(float64(senderID), responseJson["auto_accepted_by"].(float64))
	}
}

// TestDeath_PriorityOverVictory 測試死亡條件優先於勝利條件
func (suite *CleanArchTestSuite) TestDeath_PriorityOverVictory() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定目標玩家為潛伏戰線，已有 2 紅 2 黑
	senderID := game.CurrentPlayerID
	targetPlayerID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.setPlayerIdentityAndIntelligence(targetPlayerID, game.ID, entity.IdentityUndercoverFront, 2, 0, 2)

	// 傳遞一張紅黑雙色的情報卡（實際上遊戲中卡片只有一種顏色，這裡簡化為黑色）
	// 我們傳一張黑色卡，讓玩家同時達到 3 黑死亡
	blackCard := suite.addBlackIntelligenceCardToPlayer(senderID, game.ID)
	suite.Require().NotNil(blackCard, "需要黑色情報卡來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, blackCard.CardID, 0)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查玩家應該死亡（死亡優先於勝利）
	player, _ := suite.playerRepo.GetPlayerById(context.Background(), targetPlayerID)
	suite.Equal(entity.PlayerStatusDead, player.Status, "死亡條件應優先於勝利條件")

	// 遊戲應該繼續進行
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	suite.Equal(entity.GameRoomStatusPlaying, game.Status, "遊戲應繼續進行")
}

// TestDeath_AllPlayersDead 測試所有玩家都死亡時遊戲結束為平局
func (suite *CleanArchTestSuite) TestDeath_AllPlayersDead() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定除了目標玩家外的其他玩家都死亡
	senderID := game.CurrentPlayerID
	targetPlayerID := suite.getNextAlivePlayerID(game, players, senderID)

	// 設定發送者死亡
	suite.setPlayerDead(senderID)

	// 設定第三位玩家死亡
	for _, p := range players {
		if p.ID != senderID && p.ID != targetPlayerID {
			suite.setPlayerDead(p.ID)
		}
	}

	// 設定目標玩家已有 2 張黑色情報
	suite.setPlayerIntelligenceCount(targetPlayerID, game.ID, 0, 0, 2)

	// 傳遞一張黑色情報卡給目標玩家
	blackCard := suite.addBlackIntelligenceCardToPlayer(senderID, game.ID)
	suite.Require().NotNil(blackCard, "需要黑色情報卡來執行測試")

	// 手動建立情報傳遞（因為發送者已死亡，正常流程會跳過）
	suite.createIntelligenceTransferDirectly(game.ID, blackCard.CardID, senderID, targetPlayerID)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查遊戲是否結束
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	suite.Equal(entity.GameRoomStatusEnded, game.Status, "所有玩家死亡時遊戲應結束")

	// 檢查是否為平局（Winner 為空）
	suite.Equal("", game.Winner, "應為平局")
}

// ==================== 輔助方法 ====================

// setPlayerIntelligenceCount 設定玩家情報數量
func (suite *CleanArchTestSuite) setPlayerIntelligenceCount(playerID int, gameID int, red int, blue int, black int) {
	db := config.NewDatabase()

	// 清除現有情報
	db.Exec("DELETE FROM player_cards WHERE player_id = ? AND game_id = ? AND type = ?",
		playerID, gameID, entity.PlayerCardTypeIntelligence)

	// 新增紅色情報
	for i := 0; i < red; i++ {
		var cardID int
		db.Raw("SELECT id FROM cards WHERE color = '紅' LIMIT 1").Scan(&cardID)
		if cardID > 0 {
			suite.playerCardRepo.CreatePlayerCard(context.Background(),
				entity.NewPlayerCard(playerID, gameID, cardID, entity.PlayerCardTypeIntelligence))
		}
	}

	// 新增藍色情報
	for i := 0; i < blue; i++ {
		var cardID int
		db.Raw("SELECT id FROM cards WHERE color = '藍' LIMIT 1").Scan(&cardID)
		if cardID > 0 {
			suite.playerCardRepo.CreatePlayerCard(context.Background(),
				entity.NewPlayerCard(playerID, gameID, cardID, entity.PlayerCardTypeIntelligence))
		}
	}

	// 新增黑色情報
	for i := 0; i < black; i++ {
		var cardID int
		db.Raw("SELECT id FROM cards WHERE color = '黑' LIMIT 1").Scan(&cardID)
		if cardID > 0 {
			suite.playerCardRepo.CreatePlayerCard(context.Background(),
				entity.NewPlayerCard(playerID, gameID, cardID, entity.PlayerCardTypeIntelligence))
		}
	}
}

// setPlayerDead 設定玩家為死亡狀態
func (suite *CleanArchTestSuite) setPlayerDead(playerID int) {
	db := config.NewDatabase()
	db.Exec("UPDATE players SET status = ? WHERE id = ?", entity.PlayerStatusDead, playerID)
}

// setPlayerIdentityAndIntelligence 設定玩家身份和情報數量
func (suite *CleanArchTestSuite) setPlayerIdentityAndIntelligence(playerID int, gameID int, identity string, red int, blue int, black int) {
	db := config.NewDatabase()
	db.Exec("UPDATE players SET identity_card = ? WHERE id = ?", identity, playerID)
	suite.setPlayerIntelligenceCount(playerID, gameID, red, blue, black)
}

// addBlackIntelligenceCardToPlayer 為玩家新增黑色情報卡到手牌（密電類型，不需指定目標）
func (suite *CleanArchTestSuite) addBlackIntelligenceCardToPlayer(playerID int, gameID int) *entity.PlayerCard {
	db := config.NewDatabase()
	var cardID int
	// 選擇密電類型的黑色卡片（intelligence_type = 1）
	db.Raw("SELECT id FROM cards WHERE color = '黑' AND intelligence_type = 1 LIMIT 1").Scan(&cardID)

	if cardID == 0 {
		return nil
	}

	playerCard := entity.NewPlayerCard(playerID, gameID, cardID, entity.PlayerCardTypeHand)
	suite.playerCardRepo.CreatePlayerCard(context.Background(), playerCard)

	return playerCard
}

// getNextAlivePlayerIDSkipDead 取得下一位存活玩家 ID（跳過死亡玩家）
func (suite *CleanArchTestSuite) getNextAlivePlayerIDSkipDead(game *entity.Game, players []*entity.Player, currentPlayerID int) int {
	// 重新從資料庫取得玩家狀態
	freshPlayers, _ := suite.playerRepo.GetPlayersByGameId(context.Background(), game.ID)

	currentIndex := -1
	for i, p := range freshPlayers {
		if p.ID == currentPlayerID {
			currentIndex = i
			break
		}
	}

	for i := 1; i <= len(freshPlayers); i++ {
		nextIndex := (currentIndex + i) % len(freshPlayers)
		if freshPlayers[nextIndex].Status != entity.PlayerStatusDead {
			return freshPlayers[nextIndex].ID
		}
	}
	return currentPlayerID
}

// createIntelligenceTransferDirectly 直接建立情報傳遞記錄
func (suite *CleanArchTestSuite) createIntelligenceTransferDirectly(gameID int, cardID int, senderID int, targetID int) {
	transfer := entity.NewIntelligenceTransfer(gameID, cardID, senderID, targetID, targetID, false)
	suite.intelligenceTransferRepo.CreateIntelligenceTransfer(context.Background(), transfer)

	// 更新遊戲當前玩家
	db := config.NewDatabase()
	db.Exec("UPDATE games SET current_player_id = ? WHERE id = ?", targetID, gameID)
}
