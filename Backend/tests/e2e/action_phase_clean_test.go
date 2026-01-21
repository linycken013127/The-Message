package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// ==================== 抽牌測試 ====================

// TestDrawCards_Success 測試成功抽牌
func (suite *CleanArchTestSuite) TestDrawCards_Success() {
	// 建立並啟動遊戲
	game, players := suite.createAndStartGameWithPlayers(3)

	// 當前玩家抽牌
	currentPlayerID := game.CurrentPlayerID
	resp := suite.drawCards(game.ID, currentPlayerID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))
	suite.Equal(float64(2), responseJson["count"].(float64))

	cardIDs := responseJson["card_ids"].([]interface{})
	suite.Equal(2, len(cardIDs))

	_ = players // 避免 unused variable
}

// TestDrawCards_NotCurrentPlayer 測試非當前玩家抽牌
func (suite *CleanArchTestSuite) TestDrawCards_NotCurrentPlayer() {
	// 建立並啟動遊戲
	game, players := suite.createAndStartGameWithPlayers(3)

	// 找一個不是當前玩家的玩家
	var otherPlayerID int
	for _, p := range players {
		if p.ID != game.CurrentPlayerID {
			otherPlayerID = p.ID
			break
		}
	}

	// 非當前玩家嘗試抽牌
	resp := suite.drawCards(game.ID, otherPlayerID)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("還沒輪到你", responseJson["message"].(string))
}

// TestDrawCards_DeckNotEnough 測試牌堆不足（只剩 1 張）
func (suite *CleanArchTestSuite) TestDrawCards_DeckNotEnough() {
	// 建立並啟動遊戲
	game, _ := suite.createAndStartGameWithPlayers(3)

	// 取得牌堆並刪除到只剩 1 張
	decks, _ := suite.deckUseCase.GetDecksByGameId(context.Background(), game.ID)
	for i := 1; i < len(decks); i++ {
		suite.deckUseCase.DeleteDeckFromGame(context.Background(), decks[i].ID)
	}

	// 當前玩家抽牌，應該只能抽 1 張
	currentPlayerID := game.CurrentPlayerID
	resp := suite.drawCards(game.ID, currentPlayerID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))
	suite.Equal(float64(1), responseJson["count"].(float64))
}

// TestDrawCards_EmptyDeck 測試牌堆為空
func (suite *CleanArchTestSuite) TestDrawCards_EmptyDeck() {
	// 建立並啟動遊戲
	game, _ := suite.createAndStartGameWithPlayers(3)

	// 清空牌堆
	decks, _ := suite.deckUseCase.GetDecksByGameId(context.Background(), game.ID)
	for _, deck := range decks {
		suite.deckUseCase.DeleteDeckFromGame(context.Background(), deck.ID)
	}

	// 當前玩家抽牌
	currentPlayerID := game.CurrentPlayerID
	resp := suite.drawCards(game.ID, currentPlayerID)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("牌堆不足", responseJson["message"].(string))
}

// ==================== 出牌測試 ====================

// TestPlayFunctionCard_Success 測試成功出功能牌
func (suite *CleanArchTestSuite) TestPlayFunctionCard_Success() {
	// 建立並啟動遊戲
	game, _ := suite.createAndStartGameWithPlayers(3)

	// 取得當前玩家的手牌
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCards(currentPlayerID)

	// 確保手牌大於 1 張才能出牌
	suite.Greater(len(handCards), 1)

	// 出一張牌
	cardID := handCards[0].CardID
	resp := suite.playFunctionCard(game.ID, currentPlayerID, cardID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))
	suite.Equal(float64(cardID), responseJson["card_id"].(float64))
}

// TestPlayFunctionCard_LastCard 測試不能打出最後一張手牌
func (suite *CleanArchTestSuite) TestPlayFunctionCard_LastCard() {
	// 建立並啟動遊戲
	game, _ := suite.createAndStartGameWithPlayers(3)

	// 取得當前玩家的手牌
	currentPlayerID := game.CurrentPlayerID
	handCards := suite.getPlayerHandCards(currentPlayerID)

	// 刪除手牌到只剩 1 張
	for i := 1; i < len(handCards); i++ {
		suite.playerCardRepo.DeletePlayerCard(context.Background(), handCards[i].ID)
	}

	// 嘗試出最後一張牌
	cardID := handCards[0].CardID
	resp := suite.playFunctionCard(game.ID, currentPlayerID, cardID)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("不能打出最後一張手牌", responseJson["message"].(string))
}

// TestPlayFunctionCard_CardNotInHand 測試出不在手牌中的卡
func (suite *CleanArchTestSuite) TestPlayFunctionCard_CardNotInHand() {
	// 建立並啟動遊戲
	game, _ := suite.createAndStartGameWithPlayers(3)

	// 使用一個不存在於手牌的卡片 ID
	invalidCardID := 9999
	currentPlayerID := game.CurrentPlayerID
	resp := suite.playFunctionCard(game.ID, currentPlayerID, invalidCardID)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("這張牌不在你的手牌中", responseJson["message"].(string))
}

// TestPlayFunctionCard_NotCurrentPlayer 測試非當前玩家出牌
func (suite *CleanArchTestSuite) TestPlayFunctionCard_NotCurrentPlayer() {
	// 建立並啟動遊戲
	game, players := suite.createAndStartGameWithPlayers(3)

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

	// 非當前玩家嘗試出牌
	resp := suite.playFunctionCard(game.ID, otherPlayer.ID, handCards[0].CardID)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("還沒輪到你", responseJson["message"].(string))
}

// ==================== 跳過測試 ====================

// TestPass_Success 測試成功跳過
func (suite *CleanArchTestSuite) TestPass_Success() {
	// 建立並啟動遊戲
	game, _ := suite.createAndStartGameWithPlayers(3)

	// 當前玩家跳過
	currentPlayerID := game.CurrentPlayerID
	resp := suite.passAction(game.ID, currentPlayerID)

	suite.Equal(http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.True(responseJson["success"].(bool))
	suite.False(responseJson["all_passed"].(bool))
}

// TestPass_AllPlayersPass 測試所有玩家都跳過後進入情報階段
func (suite *CleanArchTestSuite) TestPass_AllPlayersPass() {
	// 建立並啟動遊戲
	game, players := suite.createAndStartGameWithPlayers(3)

	// 讓所有玩家依序跳過
	for i := 0; i < len(players); i++ {
		// 重新取得遊戲狀態
		game, _ = suite.gameRepo.GetGameWithPlayers(context.Background(), game.ID)
		currentPlayerID := game.CurrentPlayerID

		resp := suite.passAction(game.ID, currentPlayerID)
		suite.Equal(http.StatusOK, resp.StatusCode)

		responseJson := suite.responseJson(resp)
		suite.True(responseJson["success"].(bool))

		if i == len(players)-1 {
			// 最後一位玩家跳過後，應該進入情報階段
			suite.True(responseJson["all_passed"].(bool))
			suite.Equal("INTELLIGENCE", responseJson["new_phase"].(string))
		}
	}

	// 確認遊戲階段已變更為情報階段
	game, _ = suite.gameRepo.GetGameById(context.Background(), game.ID)
	suite.Equal(entity.GamePhaseIntelligence, game.Phase)
}

// TestPass_NotCurrentPlayer 測試非當前玩家跳過
func (suite *CleanArchTestSuite) TestPass_NotCurrentPlayer() {
	// 建立並啟動遊戲
	game, players := suite.createAndStartGameWithPlayers(3)

	// 找一個不是當前玩家的玩家
	var otherPlayerID int
	for _, p := range players {
		if p.ID != game.CurrentPlayerID {
			otherPlayerID = p.ID
			break
		}
	}

	// 非當前玩家嘗試跳過
	resp := suite.passAction(game.ID, otherPlayerID)

	suite.Equal(http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	suite.Equal("還沒輪到你", responseJson["message"].(string))
}

// ==================== 輔助方法 ====================

// createAndStartGameWithPlayers 建立並啟動遊戲，返回遊戲和玩家
func (suite *CleanArchTestSuite) createAndStartGameWithPlayers(playerCount int) (*entity.Game, []*entity.Player) {
	// 建立帳號
	accounts := make([]*entity.Account, playerCount)
	for i := 0; i < playerCount; i++ {
		account := &entity.Account{Name: fmt.Sprintf("player%d", i+1)}
		account, _ = suite.accountRepo.CreateAccount(context.Background(), account)
		accounts[i] = account
	}

	// 房主建立遊戲房
	game, _ := suite.gameRoomUseCase.CreateGameRoom(context.Background(), accounts[0].ID, playerCount)

	// 其他玩家加入
	for i := 1; i < playerCount; i++ {
		suite.gameRoomUseCase.JoinGameRoom(context.Background(), game.ID, accounts[i].ID)
	}

	// 房主開始遊戲
	game, _ = suite.gameRoomUseCase.StartGame(context.Background(), game.ID, accounts[0].ID)

	// 取得玩家列表
	players, _ := suite.playerUseCase.GetPlayersByGameId(context.Background(), game.ID)

	return game, players
}

// getPlayerHandCards 取得玩家手牌
func (suite *CleanArchTestSuite) getPlayerHandCards(playerID int) []entity.PlayerCard {
	player, _ := suite.playerUseCase.GetPlayerWithPlayerCards(context.Background(), playerID)
	handCards := make([]entity.PlayerCard, 0)
	for _, card := range player.PlayerCards {
		if card.Type == entity.PlayerCardTypeHand {
			handCards = append(handCards, card)
		}
	}
	return handCards
}

// drawCards 抽牌 API 呼叫
func (suite *CleanArchTestSuite) drawCards(gameID int, playerID int) *http.Response {
	body := map[string]interface{}{
		"player_id": playerID,
	}
	jsonBody, _ := json.Marshal(body)
	return suite.requestJson(fmt.Sprintf("/api/v1/games/%d/actions/draw", gameID), jsonBody, "POST")
}

// playFunctionCard 出功能牌 API 呼叫
func (suite *CleanArchTestSuite) playFunctionCard(gameID int, playerID int, cardID int) *http.Response {
	body := map[string]interface{}{
		"player_id": playerID,
		"card_id":   cardID,
	}
	jsonBody, _ := json.Marshal(body)
	return suite.requestJson(fmt.Sprintf("/api/v1/games/%d/actions/play-card", gameID), jsonBody, "POST")
}

// passAction 跳過 API 呼叫
func (suite *CleanArchTestSuite) passAction(gameID int, playerID int) *http.Response {
	body := map[string]interface{}{
		"player_id": playerID,
	}
	jsonBody, _ := json.Marshal(body)
	return suite.requestJson(fmt.Sprintf("/api/v1/games/%d/actions/pass", gameID), jsonBody, "POST")
}
