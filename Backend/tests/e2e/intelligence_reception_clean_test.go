package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// ==================== 接收情報測試 ====================

// TestAcceptIntelligence_Success 測試玩家成功接收情報
func (suite *CleanArchTestSuite) TestAcceptIntelligence_Success() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 當前玩家傳遞情報
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(currentPlayerID, entity.IntelligenceTypeSecretTelegram)
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(currentPlayerID, game.ID, entity.IntelligenceTypeSecretTelegram)
	}
	suite.Require().NotEmpty(handCards, "需要密電類型的卡片來執行測試")

	cardID := handCards[0].CardID
	suite.passIntelligenceCard(game.ID, currentPlayerID, cardID, 0)

	// 取得目標玩家（下一位玩家）
	targetPlayerID := suite.getNextAlivePlayerID(game, players, currentPlayerID)

	// 目標玩家接收情報
	resp := suite.acceptIntelligence(game.ID, targetPlayerID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))
	suite.Equal("情報已接收", responseJson["message"].(string))

	// 檢查遊戲狀態 - 應進入行動階段
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	suite.Equal(entity.GamePhaseAction, game.Phase)
	suite.Equal(targetPlayerID, game.CurrentPlayerID)

	// 檢查情報已被接收（轉為情報牌）
	playerCards := suite.getPlayerIntelligenceCards(targetPlayerID)
	suite.NotEmpty(playerCards, "玩家應該有情報牌")

	// 檢查沒有正在傳遞的情報
	transfer, _ := suite.intelligenceTransferRepo.GetActiveTransferByGameID(context.Background(), game.ID)
	suite.Nil(transfer, "不應有正在傳遞的情報")
}

// TestAcceptIntelligence_NotTargetPlayer 測試非目標玩家嘗試接收情報
func (suite *CleanArchTestSuite) TestAcceptIntelligence_NotTargetPlayer() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 當前玩家傳遞情報
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(currentPlayerID, entity.IntelligenceTypeSecretTelegram)
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(currentPlayerID, game.ID, entity.IntelligenceTypeSecretTelegram)
	}
	suite.Require().NotEmpty(handCards, "需要密電類型的卡片來執行測試")

	cardID := handCards[0].CardID
	suite.passIntelligenceCard(game.ID, currentPlayerID, cardID, 0)

	// 找一個不是目標玩家的玩家
	targetPlayerID := suite.getNextAlivePlayerID(game, players, currentPlayerID)
	var notTargetPlayerID int
	for _, p := range players {
		if p.ID != targetPlayerID && p.ID != currentPlayerID {
			notTargetPlayerID = p.ID
			break
		}
	}

	// 非目標玩家嘗試接收情報
	resp := suite.acceptIntelligence(game.ID, notTargetPlayerID)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("你不是當前情報的接收者", responseJson["message"].(string))
}

// ==================== 拒絕情報測試 ====================

// TestRejectIntelligence_Secret 測試玩家拒絕接收密電情報
func (suite *CleanArchTestSuite) TestRejectIntelligence_Secret() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 當前玩家傳遞密電情報
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(currentPlayerID, entity.IntelligenceTypeSecretTelegram)
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(currentPlayerID, game.ID, entity.IntelligenceTypeSecretTelegram)
	}
	suite.Require().NotEmpty(handCards, "需要密電類型的卡片來執行測試")

	cardID := handCards[0].CardID
	suite.passIntelligenceCard(game.ID, currentPlayerID, cardID, 0)

	// 取得目標玩家（下一位玩家）
	targetPlayerID := suite.getNextAlivePlayerID(game, players, currentPlayerID)

	// 目標玩家拒絕接收情報
	resp := suite.rejectIntelligence(game.ID, targetPlayerID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))
	suite.Equal("情報已拒絕", responseJson["message"].(string))

	// 檢查情報傳遞給下一位玩家
	transfer := responseJson["transfer"].(map[string]interface{})
	nextTargetID := suite.getNextAlivePlayerID(game, players, targetPlayerID)
	suite.Equal(float64(nextTargetID), transfer["current_target_player_id"].(float64))
}

// TestRejectIntelligence_Document 測試玩家拒絕接收文件情報
func (suite *CleanArchTestSuite) TestRejectIntelligence_Document() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 當前玩家傳遞文件情報
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(currentPlayerID, entity.IntelligenceTypeDocument)
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(currentPlayerID, game.ID, entity.IntelligenceTypeDocument)
	}
	suite.Require().NotEmpty(handCards, "需要文件類型的卡片來執行測試")

	cardID := handCards[0].CardID
	suite.passIntelligenceCard(game.ID, currentPlayerID, cardID, 0)

	// 取得目標玩家（下一位玩家）
	targetPlayerID := suite.getNextAlivePlayerID(game, players, currentPlayerID)

	// 目標玩家拒絕接收情報
	resp := suite.rejectIntelligence(game.ID, targetPlayerID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))

	// 檢查情報傳遞給下一位玩家
	transfer := responseJson["transfer"].(map[string]interface{})
	nextTargetID := suite.getNextAlivePlayerID(game, players, targetPlayerID)
	suite.Equal(float64(nextTargetID), transfer["current_target_player_id"].(float64))
}

// TestRejectIntelligence_BackToSender 測試密電情報回到發送者自動接收
func (suite *CleanArchTestSuite) TestRejectIntelligence_BackToSender() {
	// 建立遊戲並進入情報階段（3 人遊戲）
	game, players := suite.createGameInIntelligencePhase(3)

	// 當前玩家傳遞密電情報
	senderID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(senderID, entity.IntelligenceTypeSecretTelegram)
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(senderID, game.ID, entity.IntelligenceTypeSecretTelegram)
	}
	suite.Require().NotEmpty(handCards, "需要密電類型的卡片來執行測試")

	cardID := handCards[0].CardID
	suite.passIntelligenceCard(game.ID, senderID, cardID, 0)

	// 第一位目標玩家拒絕
	target1ID := suite.getNextAlivePlayerID(game, players, senderID)
	suite.rejectIntelligence(game.ID, target1ID)

	// 第二位目標玩家拒絕（此時應該回到發送者，自動接收）
	target2ID := suite.getNextAlivePlayerID(game, players, target1ID)
	resp := suite.rejectIntelligence(game.ID, target2ID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))
	// 自動接收時，訊息應該表示情報已被發送者自動接收
	suite.Equal("情報已被發送者自動接收", responseJson["message"].(string))

	// 檢查發送者有情報牌
	senderIntel := suite.getPlayerIntelligenceCards(senderID)
	suite.NotEmpty(senderIntel, "發送者應該有情報牌")

	// 檢查遊戲狀態 - 應進入行動階段
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	suite.Equal(entity.GamePhaseAction, game.Phase)

	// 檢查沒有正在傳遞的情報
	transfer, _ := suite.intelligenceTransferRepo.GetActiveTransferByGameID(context.Background(), game.ID)
	suite.Nil(transfer, "不應有正在傳遞的情報")
}

// TestRejectIntelligence_Direct 測試直達情報被拒絕後自動歸屬發送者
func (suite *CleanArchTestSuite) TestRejectIntelligence_Direct() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 當前玩家傳遞直達情報
	senderID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(senderID, entity.IntelligenceTypeDirect)
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(senderID, game.ID, entity.IntelligenceTypeDirect)
	}
	suite.Require().NotEmpty(handCards, "需要直達類型的卡片來執行測試")

	// 找一個目標玩家
	var targetPlayerID int
	for _, p := range players {
		if p.ID != senderID {
			targetPlayerID = p.ID
			break
		}
	}

	cardID := handCards[0].CardID
	suite.passIntelligenceCard(game.ID, senderID, cardID, targetPlayerID)

	// 目標玩家拒絕接收直達情報
	resp := suite.rejectIntelligence(game.ID, targetPlayerID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))
	// 直達情報被拒絕後，自動歸屬發送者
	suite.Equal("情報已被發送者自動接收", responseJson["message"].(string))

	// 檢查發送者有情報牌
	senderIntel := suite.getPlayerIntelligenceCards(senderID)
	suite.NotEmpty(senderIntel, "發送者應該有情報牌")

	// 檢查遊戲狀態 - 應進入行動階段
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	suite.Equal(entity.GamePhaseAction, game.Phase)

	// 檢查沒有正在傳遞的情報
	transfer, _ := suite.intelligenceTransferRepo.GetActiveTransferByGameID(context.Background(), game.ID)
	suite.Nil(transfer, "不應有正在傳遞的情報")
}

// TestRejectIntelligence_NotTargetPlayer 測試非目標玩家嘗試拒絕情報
func (suite *CleanArchTestSuite) TestRejectIntelligence_NotTargetPlayer() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 當前玩家傳遞情報
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(currentPlayerID, entity.IntelligenceTypeSecretTelegram)
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(currentPlayerID, game.ID, entity.IntelligenceTypeSecretTelegram)
	}
	suite.Require().NotEmpty(handCards, "需要密電類型的卡片來執行測試")

	cardID := handCards[0].CardID
	suite.passIntelligenceCard(game.ID, currentPlayerID, cardID, 0)

	// 找一個不是目標玩家的玩家
	targetPlayerID := suite.getNextAlivePlayerID(game, players, currentPlayerID)
	var notTargetPlayerID int
	for _, p := range players {
		if p.ID != targetPlayerID && p.ID != currentPlayerID {
			notTargetPlayerID = p.ID
			break
		}
	}

	// 非目標玩家嘗試拒絕情報
	resp := suite.rejectIntelligence(game.ID, notTargetPlayerID)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("你不是當前情報的接收者", responseJson["message"].(string))
}

// TestAcceptIntelligence_NoTransferInProgress 測試沒有情報傳遞時嘗試接收
func (suite *CleanArchTestSuite) TestAcceptIntelligence_NoTransferInProgress() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 不傳遞任何情報，直接嘗試接收
	resp := suite.acceptIntelligence(game.ID, players[0].ID)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("目前沒有情報傳遞中", responseJson["message"].(string))
}

// ==================== 輔助方法 ====================

// getPlayerIntelligenceCards 取得玩家的情報牌
func (suite *CleanArchTestSuite) getPlayerIntelligenceCards(playerID int) []entity.PlayerCard {
	player, _ := suite.playerUseCase.GetPlayerWithPlayerCards(context.Background(), playerID)
	intelligenceCards := make([]entity.PlayerCard, 0)
	for _, card := range player.PlayerCards {
		if card.Type == entity.PlayerCardTypeIntelligence {
			intelligenceCards = append(intelligenceCards, card)
		}
	}
	return intelligenceCards
}

// acceptIntelligence 接收情報 API 呼叫
func (suite *CleanArchTestSuite) acceptIntelligence(gameID int, playerID int) *http.Response {
	body := map[string]interface{}{
		"player_id": playerID,
	}
	jsonBody, _ := json.Marshal(body)
	return suite.requestJson(fmt.Sprintf("/api/v1/games/%d/intelligence/accept", gameID), jsonBody, "POST")
}

// rejectIntelligence 拒絕情報 API 呼叫
func (suite *CleanArchTestSuite) rejectIntelligence(gameID int, playerID int) *http.Response {
	body := map[string]interface{}{
		"player_id": playerID,
	}
	jsonBody, _ := json.Marshal(body)
	return suite.requestJson(fmt.Sprintf("/api/v1/games/%d/intelligence/reject", gameID), jsonBody, "POST")
}
