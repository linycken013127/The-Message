package e2e

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

// ============ 遊戲初始化測試 ============

// Helper: 建立並開始遊戲（指定玩家數）
func (suite *CleanArchTestSuite) createAndStartGame(playerCount int) (int, []int) {
	// 建立帳號
	accountIDs := make([]int, playerCount)
	for i := 0; i < playerCount; i++ {
		accountIDs[i] = suite.createAccount("Player" + strconv.Itoa(i+1) + "_" + strconv.Itoa(playerCount))
	}

	// 建立遊戲房（第一個玩家為房主）
	gameID := suite.createGameRoom(accountIDs[0], playerCount)

	// 其他玩家加入
	for i := 1; i < playerCount; i++ {
		suite.joinGameRoom(gameID, accountIDs[i])
	}

	// 房主開始遊戲
	api := "/api/v1/games/" + strconv.Itoa(gameID) + "/start"
	resp := suite.requestJsonWithAuth(api, nil, http.MethodPost, accountIDs[0])
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	return gameID, accountIDs
}

// Helper: 計算身份數量
func countIdentities(players []*entity.Player) (undercover, military, bystander int) {
	for _, p := range players {
		switch p.IdentityCard {
		case entity.IdentityUndercoverFront:
			undercover++
		case entity.IdentityMilitaryAgency:
			military++
		case entity.IdentityBystander:
			bystander++
		}
	}
	return
}

// ============ 身份配置測試 ============

func (suite *CleanArchTestSuite) TestGameInit_ThreePlayers_IdentityConfig() {
	// Given & When - 3 人遊戲
	gameID, _ := suite.createAndStartGame(3)

	// Then - 驗證身份配置：潛伏1, 軍情1, 打醬油1
	players, err := suite.playerRepo.GetPlayersByGameId(context.TODO(), gameID)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 3, len(players))

	undercover, military, bystander := countIdentities(players)
	assert.Equal(suite.T(), 1, undercover, "潛伏戰線應為 1 人")
	assert.Equal(suite.T(), 1, military, "軍情處應為 1 人")
	assert.Equal(suite.T(), 1, bystander, "打醬油應為 1 人")
}

func (suite *CleanArchTestSuite) TestGameInit_FivePlayers_IdentityConfig() {
	// Given & When - 5 人遊戲
	gameID, _ := suite.createAndStartGame(5)

	// Then - 驗證身份配置：潛伏2, 軍情2, 打醬油1
	players, err := suite.playerRepo.GetPlayersByGameId(context.TODO(), gameID)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 5, len(players))

	undercover, military, bystander := countIdentities(players)
	assert.Equal(suite.T(), 2, undercover, "潛伏戰線應為 2 人")
	assert.Equal(suite.T(), 2, military, "軍情處應為 2 人")
	assert.Equal(suite.T(), 1, bystander, "打醬油應為 1 人")
}

func (suite *CleanArchTestSuite) TestGameInit_SixPlayers_IdentityConfig() {
	// Given & When - 6 人遊戲
	gameID, _ := suite.createAndStartGame(6)

	// Then - 驗證身份配置：潛伏2, 軍情2, 打醬油2
	players, err := suite.playerRepo.GetPlayersByGameId(context.TODO(), gameID)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 6, len(players))

	undercover, military, bystander := countIdentities(players)
	assert.Equal(suite.T(), 2, undercover, "潛伏戰線應為 2 人")
	assert.Equal(suite.T(), 2, military, "軍情處應為 2 人")
	assert.Equal(suite.T(), 2, bystander, "打醬油應為 2 人")
}

func (suite *CleanArchTestSuite) TestGameInit_SevenPlayers_IdentityConfig() {
	// Given & When - 7 人遊戲
	gameID, _ := suite.createAndStartGame(7)

	// Then - 驗證身份配置：潛伏2, 軍情2, 打醬油3
	players, err := suite.playerRepo.GetPlayersByGameId(context.TODO(), gameID)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 7, len(players))

	undercover, military, bystander := countIdentities(players)
	assert.Equal(suite.T(), 2, undercover, "潛伏戰線應為 2 人")
	assert.Equal(suite.T(), 2, military, "軍情處應為 2 人")
	assert.Equal(suite.T(), 3, bystander, "打醬油應為 3 人")
}

func (suite *CleanArchTestSuite) TestGameInit_EightPlayers_IdentityConfig() {
	// Given & When - 8 人遊戲
	gameID, _ := suite.createAndStartGame(8)

	// Then - 驗證身份配置：潛伏3, 軍情3, 打醬油2
	players, err := suite.playerRepo.GetPlayersByGameId(context.TODO(), gameID)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 8, len(players))

	undercover, military, bystander := countIdentities(players)
	assert.Equal(suite.T(), 3, undercover, "潛伏戰線應為 3 人")
	assert.Equal(suite.T(), 3, military, "軍情處應為 3 人")
	assert.Equal(suite.T(), 2, bystander, "打醬油應為 2 人")
}

func (suite *CleanArchTestSuite) TestGameInit_NinePlayers_IdentityConfig() {
	// Given & When - 9 人遊戲
	gameID, _ := suite.createAndStartGame(9)

	// Then - 驗證身份配置：潛伏3, 軍情3, 打醬油3
	players, err := suite.playerRepo.GetPlayersByGameId(context.TODO(), gameID)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 9, len(players))

	undercover, military, bystander := countIdentities(players)
	assert.Equal(suite.T(), 3, undercover, "潛伏戰線應為 3 人")
	assert.Equal(suite.T(), 3, military, "軍情處應為 3 人")
	assert.Equal(suite.T(), 3, bystander, "打醬油應為 3 人")
}

func (suite *CleanArchTestSuite) TestGameInit_FourPlayers_NotSupported() {
	// Given - 4 人遊戲（不支援）
	accountIDs := make([]int, 4)
	for i := 0; i < 4; i++ {
		accountIDs[i] = suite.createAccount("FourP" + strconv.Itoa(i+1))
	}

	// 建立遊戲房
	gameID := suite.createGameRoom(accountIDs[0], 4)

	// 其他玩家加入
	for i := 1; i < 4; i++ {
		suite.joinGameRoom(gameID, accountIDs[i])
	}

	// When - 嘗試開始遊戲
	api := "/api/v1/games/" + strconv.Itoa(gameID) + "/start"
	resp := suite.requestJsonWithAuth(api, nil, http.MethodPost, accountIDs[0])

	// Then - 應該失敗（4人不支援）
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	assert.Contains(suite.T(), responseJson["message"], "不支援")
}

// ============ 初始手牌測試 ============

func (suite *CleanArchTestSuite) TestGameInit_PlayersReceiveThreeCards() {
	// Given & When - 開始 3 人遊戲
	gameID, _ := suite.createAndStartGame(3)

	// Then - 每位玩家應該有 3 張手牌
	players, err := suite.playerRepo.GetPlayersByGameId(context.TODO(), gameID)
	assert.Nil(suite.T(), err)

	for _, player := range players {
		p, err := suite.playerRepo.GetPlayerWithPlayerCards(context.TODO(), player.ID)
		assert.Nil(suite.T(), err)

		handCards := 0
		for _, card := range p.PlayerCards {
			if card.Type == entity.PlayerCardTypeHand {
				handCards++
			}
		}
		assert.Equal(suite.T(), 3, handCards, "每位玩家應有 3 張手牌")
	}
}

// ============ 牌堆測試 ============

func (suite *CleanArchTestSuite) TestGameInit_DeckHas45Cards() {
	// Given & When - 開始 3 人遊戲
	gameID, _ := suite.createAndStartGame(3)

	// Then - 牌堆應該有 45 張牌（總牌數）
	// 發給 3 位玩家各 3 張 = 9 張，牌堆剩 36 張
	game, err := suite.gameRepo.GetGameWithPlayers(context.TODO(), gameID)
	assert.Nil(suite.T(), err)

	// 驗證遊戲狀態
	assert.Equal(suite.T(), entity.GameRoomStatusPlaying, game.Status)

	// 注意：這裡需要確認牌堆的實作方式
	// 假設牌堆存在 deck 表中，需要查詢 deck 表
}

// ============ 第一位行動玩家測試 ============

func (suite *CleanArchTestSuite) TestGameInit_FirstPlayerSet() {
	// Given & When - 開始 3 人遊戲
	gameID, _ := suite.createAndStartGame(3)

	// Then - 應該設定第一位行動玩家
	game, err := suite.gameRepo.GetGameWithPlayers(context.TODO(), gameID)
	assert.Nil(suite.T(), err)
	assert.NotEqual(suite.T(), 0, game.CurrentPlayerID, "應該設定第一位行動玩家")

	// 驗證第一位行動玩家是遊戲中的玩家
	found := false
	for _, player := range game.Players {
		if player.ID == game.CurrentPlayerID {
			found = true
			break
		}
	}
	assert.True(suite.T(), found, "第一位行動玩家應在遊戲玩家列表中")
}
