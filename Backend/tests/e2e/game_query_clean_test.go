package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/infrastructure/config"
)

// ==================== 遊戲查詢測試 ====================

// TestGameQuery_Success 測試成功查詢遊戲資訊
func (suite *CleanArchTestSuite) TestGameQuery_Success() {
	// 建立遊戲並進入行動階段
	game, _ := suite.createGameInActionPhase(3)

	// 查詢遊戲資訊
	resp := suite.getGameInfo(game.ID, game.HostAccountID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 解析回應
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// 驗證基本欄位
	suite.Equal(float64(game.ID), result["gameId"])
	suite.Equal(entity.GameRoomStatusPlaying, result["status"])
	suite.NotNil(result["currentPlayerId"])
	suite.Equal(entity.GamePhaseAction, result["phase"])
}

// TestGameQuery_WithPersonalInfo 測試查詢結果包含玩家個人資訊
func (suite *CleanArchTestSuite) TestGameQuery_WithPersonalInfo() {
	// 建立遊戲並進入行動階段
	game, players := suite.createGameInActionPhase(3)

	// 設定玩家身份和情報
	targetPlayer := players[0]
	suite.setPlayerIdentityAndIntelligence(targetPlayer.ID, game.ID, entity.IdentityUndercoverFront, 1, 0, 1)

	// 取得該玩家的 accountID
	gamePlayers, _ := suite.gamePlayerRepo.GetGamePlayersByGameID(context.Background(), game.ID)
	var accountID int
	for _, gp := range gamePlayers {
		if gp.JoinOrder == 1 { // 第一位玩家
			accountID = gp.AccountID
			break
		}
	}

	// 查詢遊戲資訊
	resp := suite.getGameInfo(game.ID, accountID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 解析回應
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// 驗證個人資訊
	suite.Equal(entity.IdentityUndercoverFront, result["myFaction"])
	suite.Equal(float64(1), result["myRedCount"])
	suite.Equal(float64(0), result["myBlueCount"])
	suite.Equal(float64(1), result["myBlackCount"])
	suite.NotNil(result["myHandCardCount"])
}

// TestGameQuery_Unauthorized 測試未認證無法查詢
func (suite *CleanArchTestSuite) TestGameQuery_Unauthorized() {
	// 建立遊戲
	game, _ := suite.createGameInActionPhase(3)

	// 不帶認證查詢遊戲資訊
	resp := suite.getGameInfoWithoutAuth(game.ID)
	suite.Equal(http.StatusUnauthorized, resp.StatusCode)
}

// TestGameQuery_WithIntelligenceInTransit 測試查詢結果包含情報傳遞狀態
func (suite *CleanArchTestSuite) TestGameQuery_WithIntelligenceInTransit() {
	// 建立遊戲並進入情報階段
	game, _ := suite.createGameInIntelligencePhase(3)

	// 為當前玩家新增一張密電卡
	senderID := game.CurrentPlayerID
	cards := suite.addIntelligenceCardToPlayer(senderID, game.ID, entity.IntelligenceTypeSecretTelegram)
	suite.Require().NotEmpty(cards, "需要情報卡來執行測試")

	// 傳遞情報
	resp := suite.passIntelligenceCard(game.ID, senderID, cards[0].CardID, 0)
	suite.Equal(http.StatusOK, resp.StatusCode, "傳遞情報應該成功")

	// 查詢遊戲資訊
	resp = suite.getGameInfo(game.ID, game.HostAccountID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 解析回應
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// 驗證有情報正在傳遞
	suite.Equal(true, result["hasIntelInTransit"])
	suite.Equal(entity.GamePhaseIntelligence, result["phase"])
}

// TestGameQuery_WithPlayersInfo 測試查詢結果包含所有玩家公開資訊
func (suite *CleanArchTestSuite) TestGameQuery_WithPlayersInfo() {
	// 建立遊戲並進入行動階段
	game, _ := suite.createGameInActionPhase(3)

	// 查詢遊戲資訊
	resp := suite.getGameInfo(game.ID, game.HostAccountID)
	suite.Equal(http.StatusOK, resp.StatusCode)

	// 解析回應
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// 驗證玩家列表
	players, ok := result["players"].([]interface{})
	suite.True(ok, "應該包含玩家列表")
	suite.Equal(3, len(players), "應該有 3 位玩家")

	// 檢查每個玩家的公開資訊
	for _, p := range players {
		playerInfo := p.(map[string]interface{})
		suite.NotNil(playerInfo["id"])
		suite.NotNil(playerInfo["name"])
		suite.NotNil(playerInfo["alive"])
		suite.NotNil(playerInfo["intelCount"])
	}
}

// TestGameQuery_NotInGame 測試查詢不存在的遊戲
func (suite *CleanArchTestSuite) TestGameQuery_NotInGame() {
	// 建立帳號
	account, _ := suite.accountUseCase.RegisterAccount(context.Background(), "TestPlayer")

	// 查詢不存在的遊戲
	resp := suite.getGameInfo(99999, account.ID)
	suite.Equal(http.StatusNotFound, resp.StatusCode)
}

// ==================== 輔助方法 ====================

// createGameInActionPhase 建立遊戲並進入行動階段
func (suite *CleanArchTestSuite) createGameInActionPhase(playerCount int) (*entity.Game, []*entity.Player) {
	ctx := context.Background()

	// 建立帳號
	accounts := make([]*entity.Account, playerCount)
	for i := 0; i < playerCount; i++ {
		account, err := suite.accountUseCase.RegisterAccount(ctx, fmt.Sprintf("Player%d", i+1))
		suite.Require().NoError(err)
		accounts[i] = account
	}

	// 建立遊戲房
	game, err := suite.gameRoomUseCase.CreateGameRoom(ctx, accounts[0].ID, 9)
	suite.Require().NoError(err)

	// 其他玩家加入
	for i := 1; i < playerCount; i++ {
		err = suite.gameRoomUseCase.JoinGameRoom(ctx, game.ID, accounts[i].ID)
		suite.Require().NoError(err)
	}

	// 開始遊戲
	game, err = suite.gameRoomUseCase.StartGame(ctx, game.ID, accounts[0].ID)
	suite.Require().NoError(err)

	// 取得所有玩家
	players, err := suite.playerRepo.GetPlayersByGameId(ctx, game.ID)
	suite.Require().NoError(err)

	return game, players
}

// getGameInfo 查詢遊戲資訊（帶認證）
func (suite *CleanArchTestSuite) getGameInfo(gameID int, accountID int) *http.Response {
	url := fmt.Sprintf("%s/api/v1/games/%d", suite.server.URL, gameID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Account-ID", fmt.Sprintf("%d", accountID))

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)

	return resp
}

// getGameInfoWithoutAuth 查詢遊戲資訊（不帶認證）
func (suite *CleanArchTestSuite) getGameInfoWithoutAuth(gameID int) *http.Response {
	url := fmt.Sprintf("%s/api/v1/games/%d", suite.server.URL, gameID)
	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)

	return resp
}

// setPlayerIdentityAndIntelligence 設定玩家身份和情報（從 death_conditions_clean_test.go 移過來的共用方法）
func (suite *CleanArchTestSuite) setPlayerIdentityAndIntelligenceForQuery(playerID int, gameID int, identity string, redCount int, blueCount int, blackCount int) {
	ctx := context.Background()

	// 更新玩家身份
	player, err := suite.playerRepo.GetPlayerById(ctx, playerID)
	suite.Require().NoError(err)
	player.IdentityCard = identity
	suite.playerRepo.UpdatePlayer(ctx, player)

	// 新增紅色情報
	db := config.NewDatabase()
	for i := 0; i < redCount; i++ {
		var cardID int
		db.Raw("SELECT id FROM cards WHERE color = '紅' AND intelligence_type = 1 LIMIT 1 OFFSET ?", i).Scan(&cardID)
		if cardID > 0 {
			pc := entity.NewPlayerCard(playerID, gameID, cardID, entity.PlayerCardTypeIntelligence)
			suite.playerCardRepo.CreatePlayerCard(ctx, pc)
		}
	}

	// 新增藍色情報
	for i := 0; i < blueCount; i++ {
		var cardID int
		db.Raw("SELECT id FROM cards WHERE color = '藍' AND intelligence_type = 1 LIMIT 1 OFFSET ?", i).Scan(&cardID)
		if cardID > 0 {
			pc := entity.NewPlayerCard(playerID, gameID, cardID, entity.PlayerCardTypeIntelligence)
			suite.playerCardRepo.CreatePlayerCard(ctx, pc)
		}
	}

	// 新增黑色情報
	for i := 0; i < blackCount; i++ {
		var cardID int
		db.Raw("SELECT id FROM cards WHERE color = '黑' AND intelligence_type = 1 LIMIT 1 OFFSET ?", i).Scan(&cardID)
		if cardID > 0 {
			pc := entity.NewPlayerCard(playerID, gameID, cardID, entity.PlayerCardTypeIntelligence)
			suite.playerCardRepo.CreatePlayerCard(ctx, pc)
		}
	}
}
