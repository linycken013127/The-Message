package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/config"
)

// ==================== 傳遞情報牌測試 ====================

// TestPassIntelligenceCard_SecretTelegram 測試傳遞密電（蓋牌向右傳）
func (suite *CleanArchTestSuite) TestPassIntelligenceCard_SecretTelegram() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 取得當前玩家（送出情報的玩家）
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(currentPlayerID, entity.IntelligenceTypeSecretTelegram)

	// 如果沒有密電卡，手動設定一張
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(currentPlayerID, game.ID, entity.IntelligenceTypeSecretTelegram)
	}
	suite.Require().NotEmpty(handCards, "需要密電類型的卡片來執行測試")

	// 傳遞密電
	cardID := handCards[0].CardID
	resp := suite.passIntelligenceCard(game.ID, currentPlayerID, cardID, 0)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))
	suite.Equal("情報牌已送出", responseJson["message"].(string))

	// 檢查傳遞狀態
	transfer := responseJson["transfer"].(map[string]interface{})
	suite.Equal(float64(game.ID), transfer["game_id"].(float64))
	suite.Equal(float64(cardID), transfer["card_id"].(float64))
	suite.Equal(float64(currentPlayerID), transfer["sender_player_id"].(float64))
	suite.False(transfer["face_up"].(bool)) // 密電是蓋牌
	suite.Equal("IN_TRANSIT", transfer["status"].(string))

	// 目標應該是下一位玩家（向右傳）
	targetPlayerID := suite.getNextAlivePlayerID(game, players, currentPlayerID)
	suite.Equal(float64(targetPlayerID), transfer["current_target_player_id"].(float64))
}

// TestPassIntelligenceCard_Document 測試傳遞文件（明牌向右傳）
func (suite *CleanArchTestSuite) TestPassIntelligenceCard_Document() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 取得當前玩家
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(currentPlayerID, entity.IntelligenceTypeDocument)

	// 如果沒有文件卡，手動設定一張
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(currentPlayerID, game.ID, entity.IntelligenceTypeDocument)
	}
	suite.Require().NotEmpty(handCards, "需要文件類型的卡片來執行測試")

	// 傳遞文件
	cardID := handCards[0].CardID
	resp := suite.passIntelligenceCard(game.ID, currentPlayerID, cardID, 0)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))

	// 檢查傳遞狀態
	transfer := responseJson["transfer"].(map[string]interface{})
	suite.True(transfer["face_up"].(bool)) // 文件是明牌

	// 目標應該是下一位玩家（向右傳）
	targetPlayerID := suite.getNextAlivePlayerID(game, players, currentPlayerID)
	suite.Equal(float64(targetPlayerID), transfer["current_target_player_id"].(float64))

	// 明牌應該顯示卡片資訊
	suite.NotNil(responseJson["card"])
}

// TestPassIntelligenceCard_Direct 測試傳遞直達（蓋牌指定玩家）
func (suite *CleanArchTestSuite) TestPassIntelligenceCard_Direct() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 取得當前玩家
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(currentPlayerID, entity.IntelligenceTypeDirect)

	// 如果沒有直達卡，手動設定一張
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(currentPlayerID, game.ID, entity.IntelligenceTypeDirect)
	}

	// 確保有可用的卡片
	suite.Require().NotEmpty(handCards, "需要直達類型的卡片來執行測試")

	// 找一個不是當前玩家的目標
	var targetPlayerID int
	for _, p := range players {
		if p.ID != currentPlayerID {
			targetPlayerID = p.ID
			break
		}
	}

	// 傳遞直達
	cardID := handCards[0].CardID
	resp := suite.passIntelligenceCard(game.ID, currentPlayerID, cardID, targetPlayerID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))

	// 檢查傳遞狀態
	transfer := responseJson["transfer"].(map[string]interface{})
	suite.False(transfer["face_up"].(bool)) // 直達是蓋牌
	suite.Equal(float64(targetPlayerID), transfer["current_target_player_id"].(float64))
}

// TestPassIntelligenceCard_DirectWithoutTarget 測試直達卡未指定目標
func (suite *CleanArchTestSuite) TestPassIntelligenceCard_DirectWithoutTarget() {
	// 建立遊戲並進入情報階段
	game, _ := suite.createGameInIntelligencePhase(3)

	// 取得當前玩家
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(currentPlayerID, entity.IntelligenceTypeDirect)

	// 如果沒有直達卡，手動設定一張
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(currentPlayerID, game.ID, entity.IntelligenceTypeDirect)
	}
	suite.Require().NotEmpty(handCards, "需要直達類型的卡片來執行測試")

	// 傳遞直達但不指定目標
	cardID := handCards[0].CardID
	resp := suite.passIntelligenceCard(game.ID, currentPlayerID, cardID, 0)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("直達情報需要指定目標玩家", responseJson["message"].(string))
}

// TestPassIntelligenceCard_DirectCannotTargetSelf 測試直達卡不能指定自己
func (suite *CleanArchTestSuite) TestPassIntelligenceCard_DirectCannotTargetSelf() {
	// 建立遊戲並進入情報階段
	game, _ := suite.createGameInIntelligencePhase(3)

	// 取得當前玩家
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(currentPlayerID, entity.IntelligenceTypeDirect)

	// 如果沒有直達卡，手動設定一張
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(currentPlayerID, game.ID, entity.IntelligenceTypeDirect)
	}
	suite.Require().NotEmpty(handCards, "需要直達類型的卡片來執行測試")

	// 傳遞直達指定自己
	cardID := handCards[0].CardID
	resp := suite.passIntelligenceCard(game.ID, currentPlayerID, cardID, currentPlayerID)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("直達情報不能指定自己", responseJson["message"].(string))
}

// TestPassIntelligenceCard_NotCurrentPlayer 測試非當前玩家傳遞情報
func (suite *CleanArchTestSuite) TestPassIntelligenceCard_NotCurrentPlayer() {
	// 建立遊戲並進入情報階段
	game, players := suite.createGameInIntelligencePhase(3)

	// 找一個不是當前玩家的玩家
	var otherPlayer *entity.Player
	for _, p := range players {
		if p.ID != game.CurrentPlayerID {
			otherPlayer = p
			break
		}
	}

	// 取得該玩家的手牌
	handCards := suite.getPlayerHandCards(otherPlayer.ID)
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(otherPlayer.ID, game.ID, entity.IntelligenceTypeSecretTelegram)
	}

	// 非當前玩家嘗試傳遞情報
	resp := suite.passIntelligenceCard(game.ID, otherPlayer.ID, handCards[0].CardID, 0)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("還沒輪到你", responseJson["message"].(string))
}

// TestPassIntelligenceCard_NotInIntelligencePhase 測試非情報階段傳遞情報
func (suite *CleanArchTestSuite) TestPassIntelligenceCard_NotInIntelligencePhase() {
	// 建立並啟動遊戲（此時在行動階段）
	game, _ := suite.createAndStartGameWithPlayers(3)

	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCards(currentPlayerID)

	// 嘗試在行動階段傳遞情報
	resp := suite.passIntelligenceCard(game.ID, currentPlayerID, handCards[0].CardID, 0)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("不在情報階段", responseJson["message"].(string))
}

// TestPassIntelligenceCard_CardNotInHand 測試傳遞不在手牌中的卡
func (suite *CleanArchTestSuite) TestPassIntelligenceCard_CardNotInHand() {
	// 建立遊戲並進入情報階段
	game, _ := suite.createGameInIntelligencePhase(3)

	currentPlayerID := game.CurrentPlayerID

	// 使用一個不存在於手牌的卡片 ID
	invalidCardID := 9999
	resp := suite.passIntelligenceCard(game.ID, currentPlayerID, invalidCardID, 0)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("這張牌不在你的手牌中", responseJson["message"].(string))
}

// ==================== 取得正在傳遞的情報測試 ====================

// TestGetActiveTransfer_HasActiveTransfer 測試有正在傳遞的情報
func (suite *CleanArchTestSuite) TestGetActiveTransfer_HasActiveTransfer() {
	// 建立遊戲並進入情報階段
	game, _ := suite.createGameInIntelligencePhase(3)

	// 當前玩家傳遞情報
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCardsWithType(currentPlayerID, entity.IntelligenceTypeSecretTelegram)
	if len(handCards) == 0 {
		handCards = suite.addIntelligenceCardToPlayer(currentPlayerID, game.ID, entity.IntelligenceTypeSecretTelegram)
	}
	suite.Require().NotEmpty(handCards, "需要密電類型的卡片來執行測試")

	suite.passIntelligenceCard(game.ID, currentPlayerID, handCards[0].CardID, 0)

	// 取得正在傳遞的情報
	resp := suite.getActiveTransfer(game.ID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["has_active_transfer"].(bool))
	suite.NotNil(responseJson["transfer"])
}

// TestGetActiveTransfer_NoActiveTransfer 測試沒有正在傳遞的情報
func (suite *CleanArchTestSuite) TestGetActiveTransfer_NoActiveTransfer() {
	// 建立遊戲並進入情報階段
	game, _ := suite.createGameInIntelligencePhase(3)

	// 不傳遞任何情報，直接查詢
	resp := suite.getActiveTransfer(game.ID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.False(responseJson["has_active_transfer"].(bool))
}

// ==================== 輔助方法 ====================

// createGameInIntelligencePhase 建立遊戲並進入情報階段
func (suite *CleanArchTestSuite) createGameInIntelligencePhase(playerCount int) (*entity.Game, []*entity.Player) {
	// 建立並啟動遊戲
	game, players := suite.createAndStartGameWithPlayers(playerCount)

	// 讓所有玩家都跳過以進入情報階段
	for i := 0; i < playerCount; i++ {
		game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
		currentPlayerID := game.CurrentPlayerID
		suite.passAction(game.ID, currentPlayerID)
	}

	// 重新取得遊戲狀態
	game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
	return game, players
}

// getPlayerHandCardsWithType 取得玩家特定類型的手牌
func (suite *CleanArchTestSuite) getPlayerHandCardsWithType(playerID int, intelligenceType int) []entity.PlayerCard {
	player, _ := suite.playerUseCase.GetPlayerWithPlayerCards(context.Background(), playerID)
	handCards := make([]entity.PlayerCard, 0)
	for _, card := range player.PlayerCards {
		if card.Type == entity.PlayerCardTypeHand && card.Card.ID != 0 && card.Card.IntelligenceType == intelligenceType {
			handCards = append(handCards, card)
		}
	}
	return handCards
}

// addIntelligenceCardToPlayer 為玩家添加特定類型的情報卡
func (suite *CleanArchTestSuite) addIntelligenceCardToPlayer(playerID int, gameID int, intelligenceType int) []entity.PlayerCard {
	// 使用新的資料庫連線查詢卡片
	db := config.NewDatabase()
	var cardID int
	db.Raw("SELECT id FROM cards WHERE intelligence_type = ? LIMIT 1", intelligenceType).Scan(&cardID)

	if cardID == 0 {
		return nil
	}

	// 建立玩家手牌
	playerCard := entity.NewPlayerCard(playerID, gameID, cardID, entity.PlayerCardTypeHand)
	suite.playerCardRepo.CreatePlayerCard(context.Background(), playerCard)

	return []entity.PlayerCard{*playerCard}
}

// getNextAlivePlayerID 取得下一位存活玩家 ID（向右）
func (suite *CleanArchTestSuite) getNextAlivePlayerID(game *entity.Game, players []*entity.Player, currentPlayerID int) int {
	currentIndex := -1
	for i, p := range players {
		if p.ID == currentPlayerID {
			currentIndex = i
			break
		}
	}

	for i := 1; i <= len(players); i++ {
		nextIndex := (currentIndex + i) % len(players)
		if players[nextIndex].Status != entity.PlayerStatusDead {
			return players[nextIndex].ID
		}
	}
	return currentPlayerID
}

// passIntelligenceCard 傳遞情報牌 API 呼叫
func (suite *CleanArchTestSuite) passIntelligenceCard(gameID int, playerID int, cardID int, targetPlayerID int) *http.Response {
	body := map[string]interface{}{
		"player_id": playerID,
		"card_id":   cardID,
	}
	if targetPlayerID > 0 {
		body["target_player_id"] = targetPlayerID
	}
	jsonBody, _ := json.Marshal(body)
	return suite.requestJson(fmt.Sprintf("/api/v1/games/%d/intelligence/pass-card", gameID), jsonBody, "POST")
}

// getActiveTransfer 取得正在傳遞的情報 API 呼叫
func (suite *CleanArchTestSuite) getActiveTransfer(gameID int) *http.Response {
	return suite.requestJson(fmt.Sprintf("/api/v1/games/%d/intelligence/active", gameID), nil, "GET")
}
