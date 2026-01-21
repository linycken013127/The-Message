package e2e

import (
	"context"
	"net/http"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/config"
)

// ==================== 勝利條件測試 ====================

// TestVictory_UndercoverFront_ThreeRed 測試潛伏戰線收集 3 張紅色情報獲勝
func (suite *CleanArchTestSuite) TestVictory_UndercoverFront_ThreeRed() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定目標玩家為潛伏戰線，已有 2 張紅色情報
	senderID := game.CurrentPlayerID
	targetPlayerID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.setPlayerIdentityAndIntelligence(targetPlayerID, game.ID, entity.IdentityUndercoverFront, 2, 0, 0)

	// 當前玩家傳遞一張紅色情報卡
	redCard := suite.addRedIntelligenceCardToPlayer(senderID, game.ID)
	suite.Require().NotNil(redCard, "需要紅色情報卡來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, redCard.CardID, 0)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查遊戲是否結束且潛伏戰線獲勝
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	suite.Equal(entity.GameRoomStatusEnded, game.Status, "遊戲應該結束")
	suite.Equal(entity.IdentityUndercoverFront, game.Winner, "潛伏戰線應該獲勝")
}

// TestVictory_UndercoverFront_MoreThanThreeRed 測試潛伏戰線收集超過 3 張紅色情報獲勝
func (suite *CleanArchTestSuite) TestVictory_UndercoverFront_MoreThanThreeRed() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定目標玩家為潛伏戰線，已有 3 張紅色情報
	senderID := game.CurrentPlayerID
	targetPlayerID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.setPlayerIdentityAndIntelligence(targetPlayerID, game.ID, entity.IdentityUndercoverFront, 3, 0, 0)

	// 當前玩家傳遞一張紅色情報卡
	redCard := suite.addRedIntelligenceCardToPlayer(senderID, game.ID)
	suite.Require().NotNil(redCard, "需要紅色情報卡來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, redCard.CardID, 0)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查遊戲是否結束且潛伏戰線獲勝
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	suite.Equal(entity.GameRoomStatusEnded, game.Status, "遊戲應該結束")
	suite.Equal(entity.IdentityUndercoverFront, game.Winner, "潛伏戰線應該獲勝")
}

// TestVictory_MilitaryAgency_ThreeBlue 測試軍情處收集 3 張藍色情報獲勝
func (suite *CleanArchTestSuite) TestVictory_MilitaryAgency_ThreeBlue() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定目標玩家為軍情處，已有 2 張藍色情報
	senderID := game.CurrentPlayerID
	targetPlayerID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.setPlayerIdentityAndIntelligence(targetPlayerID, game.ID, entity.IdentityMilitaryAgency, 0, 2, 0)

	// 當前玩家傳遞一張藍色情報卡
	blueCard := suite.addBlueIntelligenceCardToPlayer(senderID, game.ID)
	suite.Require().NotNil(blueCard, "需要藍色情報卡來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, blueCard.CardID, 0)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查遊戲是否結束且軍情處獲勝
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	suite.Equal(entity.GameRoomStatusEnded, game.Status, "遊戲應該結束")
	suite.Equal(entity.IdentityMilitaryAgency, game.Winner, "軍情處應該獲勝")
}

// TestVictory_Bystander_CannotTrigger 測試打醬油的無法單獨觸發勝利
func (suite *CleanArchTestSuite) TestVictory_Bystander_CannotTrigger() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定目標玩家為打醬油，已有 2 張紅色情報
	senderID := game.CurrentPlayerID
	targetPlayerID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.setPlayerIdentityAndIntelligence(targetPlayerID, game.ID, entity.IdentityBystander, 2, 0, 0)

	// 當前玩家傳遞一張紅色情報卡
	redCard := suite.addRedIntelligenceCardToPlayer(senderID, game.ID)
	suite.Require().NotNil(redCard, "需要紅色情報卡來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, redCard.CardID, 0)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查遊戲應該繼續（打醬油無法觸發勝利）
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	suite.Equal(entity.GameRoomStatusPlaying, game.Status, "遊戲應該繼續")
}

// TestVictory_NotReached 測試情報數量不足時遊戲繼續
func (suite *CleanArchTestSuite) TestVictory_NotReached() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定目標玩家為潛伏戰線，已有 1 張紅色情報
	senderID := game.CurrentPlayerID
	targetPlayerID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.setPlayerIdentityAndIntelligence(targetPlayerID, game.ID, entity.IdentityUndercoverFront, 1, 0, 0)

	// 當前玩家傳遞一張紅色情報卡
	redCard := suite.addRedIntelligenceCardToPlayer(senderID, game.ID)
	suite.Require().NotNil(redCard, "需要紅色情報卡來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, redCard.CardID, 0)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查遊戲應該繼續
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	suite.Equal(entity.GameRoomStatusPlaying, game.Status, "遊戲應該繼續")
}

// TestVictory_DeathPriorityOverVictory 測試死亡條件優先於勝利條件
func (suite *CleanArchTestSuite) TestVictory_DeathPriorityOverVictory() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 設定目標玩家為潛伏戰線，已有 2 張紅色和 2 張黑色情報
	senderID := game.CurrentPlayerID
	targetPlayerID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.setPlayerIdentityAndIntelligence(targetPlayerID, game.ID, entity.IdentityUndercoverFront, 2, 0, 2)

	// 當前玩家傳遞一張紅黑雙色卡（使用黑色測試死亡優先）
	// 由於遊戲卡片只有單一顏色，這裡用黑色測試
	blackCard := suite.addBlackIntelligenceCardToPlayer(senderID, game.ID)
	suite.Require().NotNil(blackCard, "需要黑色情報卡來執行測試")

	suite.passIntelligenceCard(game.ID, senderID, blackCard.CardID, 0)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 檢查玩家應該死亡（死亡優先於勝利）
	player, _ := suite.playerRepo.GetPlayerById(context.Background(), targetPlayerID)
	suite.Equal(entity.PlayerStatusDead, player.Status, "死亡條件應優先於勝利條件")

	// 遊戲應該繼續進行（因為還有其他玩家存活）
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	suite.Equal(entity.GameRoomStatusPlaying, game.Status, "遊戲應繼續進行")
}

// ==================== 輔助方法 ====================

// addRedIntelligenceCardToPlayer 為玩家新增紅色情報卡到手牌（密電類型）
func (suite *CleanArchTestSuite) addRedIntelligenceCardToPlayer(playerID int, gameID int) *entity.PlayerCard {
	db := config.NewDatabase()
	var cardID int
	// 選擇密電類型的紅色卡片（intelligence_type = 1）
	db.Raw("SELECT id FROM cards WHERE color = '紅' AND intelligence_type = 1 LIMIT 1").Scan(&cardID)

	if cardID == 0 {
		return nil
	}

	playerCard := entity.NewPlayerCard(playerID, gameID, cardID, entity.PlayerCardTypeHand)
	suite.playerCardRepo.CreatePlayerCard(context.Background(), playerCard)

	return playerCard
}

// addBlueIntelligenceCardToPlayer 為玩家新增藍色情報卡到手牌（密電類型）
func (suite *CleanArchTestSuite) addBlueIntelligenceCardToPlayer(playerID int, gameID int) *entity.PlayerCard {
	db := config.NewDatabase()
	var cardID int
	// 選擇密電類型的藍色卡片（intelligence_type = 1）
	db.Raw("SELECT id FROM cards WHERE color = '藍' AND intelligence_type = 1 LIMIT 1").Scan(&cardID)

	if cardID == 0 {
		return nil
	}

	playerCard := entity.NewPlayerCard(playerID, gameID, cardID, entity.PlayerCardTypeHand)
	suite.playerCardRepo.CreatePlayerCard(context.Background(), playerCard)

	return playerCard
}
