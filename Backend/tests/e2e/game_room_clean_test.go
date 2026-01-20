package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/stretchr/testify/assert"
)

// CreateGameRoomRequest 建立遊戲房請求
type CreateGameRoomRequest struct {
	MaxPlayers int `json:"maxPlayers"`
}

// JoinGameRoomRequest 加入遊戲房請求
type JoinGameRoomRequest struct{}

// StartGameRoomRequest 開始遊戲請求
type StartGameRoomRequest struct{}

// Helper: 建立帳號並回傳 ID
func (suite *CleanArchTestSuite) createAccount(name string) int {
	req := RegisterPlayerRequest{PlayerName: name}
	jsonBody, _ := json.Marshal(req)
	resp := suite.requestJson("/api/v1/players", jsonBody, http.MethodPost)
	responseJson := suite.responseJson(resp)
	return int(responseJson["playerId"].(float64))
}

// Helper: 建立遊戲房
func (suite *CleanArchTestSuite) createGameRoom(accountID int, maxPlayers int) int {
	api := "/api/v1/games"
	req := CreateGameRoomRequest{MaxPlayers: maxPlayers}
	jsonBody, _ := json.Marshal(req)
	resp := suite.requestJsonWithAuth(api, jsonBody, http.MethodPost, accountID)
	responseJson := suite.responseJson(resp)
	return int(responseJson["gameId"].(float64))
}

// Helper: 加入遊戲房
func (suite *CleanArchTestSuite) joinGameRoom(gameID int, accountID int) *http.Response {
	api := "/api/v1/games/" + strconv.Itoa(gameID) + "/join"
	req := JoinGameRoomRequest{}
	jsonBody, _ := json.Marshal(req)
	return suite.requestJsonWithAuth(api, jsonBody, http.MethodPost, accountID)
}

// Helper: 發送帶認證的請求
func (suite *CleanArchTestSuite) requestJsonWithAuth(api string, jsonBody []byte, method string, accountID int) *http.Response {
	req, err := http.NewRequest(method, suite.server.URL+api, bytes.NewBuffer(jsonBody))
	if err != nil {
		suite.T().Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Account-ID", strconv.Itoa(accountID))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		suite.T().Fatalf("Failed to send request: %v", err)
	}
	return resp
}

// ============ 建立遊戲房測試 ============

func (suite *CleanArchTestSuite) TestCreateGameRoom_Success() {
	// Given
	accountID := suite.createAccount("Host1")

	// When
	api := "/api/v1/games"
	req := CreateGameRoomRequest{MaxPlayers: 9}
	jsonBody, _ := json.Marshal(req)
	resp := suite.requestJsonWithAuth(api, jsonBody, http.MethodPost, accountID)

	// Then
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	assert.NotNil(suite.T(), responseJson["gameId"])
	assert.Equal(suite.T(), "WAITING", responseJson["status"])
	assert.Equal(suite.T(), float64(accountID), responseJson["hostPlayerId"])

	// 驗證資料庫
	gameID := int(responseJson["gameId"].(float64))
	game, err := suite.gameRepo.GetGameById(context.TODO(), gameID)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "WAITING", game.Status)
	assert.Equal(suite.T(), accountID, game.HostAccountID)
	assert.Equal(suite.T(), 1, game.CurrentPlayers)
}

func (suite *CleanArchTestSuite) TestCreateGameRoom_Unauthorized() {
	// When - 不帶認證
	api := "/api/v1/games"
	req := CreateGameRoomRequest{MaxPlayers: 9}
	jsonBody, _ := json.Marshal(req)
	resp := suite.requestJson(api, jsonBody, http.MethodPost)

	// Then - 應該失敗 (現有 API 會建立遊戲，但新 API 應要求認證)
	// 注意：這個測試可能需要根據實際實作調整
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

// ============ 加入遊戲房測試 ============

func (suite *CleanArchTestSuite) TestJoinGameRoom_Success() {
	// Given
	hostID := suite.createAccount("Host2")
	playerID := suite.createAccount("Player2")
	gameID := suite.createGameRoom(hostID, 9)

	// When
	resp := suite.joinGameRoom(gameID, playerID)

	// Then
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	assert.Equal(suite.T(), "加入成功", responseJson["message"])

	// 驗證資料庫
	game, _ := suite.gameRepo.GetGameById(context.TODO(), gameID)
	assert.Equal(suite.T(), 2, game.CurrentPlayers)
}

func (suite *CleanArchTestSuite) TestJoinGameRoom_AlreadyJoined() {
	// Given
	hostID := suite.createAccount("Host3")
	playerID := suite.createAccount("Player3")
	gameID := suite.createGameRoom(hostID, 9)
	suite.joinGameRoom(gameID, playerID)

	// When - 嘗試重複加入
	resp := suite.joinGameRoom(gameID, playerID)

	// Then
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestJoinGameRoom_GameFull() {
	// Given - 建立一個最多 3 人的遊戲房
	hostID := suite.createAccount("Host4")
	player1ID := suite.createAccount("P4_1")
	player2ID := suite.createAccount("P4_2")
	player3ID := suite.createAccount("P4_3")

	gameID := suite.createGameRoom(hostID, 3) // 最多 3 人
	suite.joinGameRoom(gameID, player1ID)
	suite.joinGameRoom(gameID, player2ID)

	// When - 第 4 位玩家嘗試加入
	resp := suite.joinGameRoom(gameID, player3ID)

	// Then
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestJoinGameRoom_GameAlreadyStarted() {
	// Given
	hostID := suite.createAccount("Host5")
	player1ID := suite.createAccount("P5_1")
	player2ID := suite.createAccount("P5_2")
	latePlayerID := suite.createAccount("Late5")

	gameID := suite.createGameRoom(hostID, 9)
	suite.joinGameRoom(gameID, player1ID)
	suite.joinGameRoom(gameID, player2ID)

	// 開始遊戲
	startAPI := "/api/v1/games/" + strconv.Itoa(gameID) + "/start"
	resp := suite.requestJsonWithAuth(startAPI, nil, http.MethodPost, hostID)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// When - 遲到的玩家嘗試加入
	resp = suite.joinGameRoom(gameID, latePlayerID)

	// Then
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

// ============ 開始遊戲測試 ============

func (suite *CleanArchTestSuite) TestStartGame_AsHost_Success() {
	// Given
	hostID := suite.createAccount("Host6")
	player1ID := suite.createAccount("P6_1")
	player2ID := suite.createAccount("P6_2")

	gameID := suite.createGameRoom(hostID, 9)
	suite.joinGameRoom(gameID, player1ID)
	suite.joinGameRoom(gameID, player2ID)

	// When
	api := "/api/v1/games/" + strconv.Itoa(gameID) + "/start"
	resp := suite.requestJsonWithAuth(api, nil, http.MethodPost, hostID)

	// Then
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	assert.Equal(suite.T(), "PLAYING", responseJson["status"])

	// 驗證資料庫
	game, _ := suite.gameRepo.GetGameById(context.TODO(), gameID)
	assert.Equal(suite.T(), "PLAYING", game.Status)
}

func (suite *CleanArchTestSuite) TestStartGame_NotHost_Fail() {
	// Given
	hostID := suite.createAccount("Host7")
	player1ID := suite.createAccount("P7_1")
	player2ID := suite.createAccount("P7_2")

	gameID := suite.createGameRoom(hostID, 9)
	suite.joinGameRoom(gameID, player1ID)
	suite.joinGameRoom(gameID, player2ID)

	// When - 非房主嘗試開始遊戲
	api := "/api/v1/games/" + strconv.Itoa(gameID) + "/start"
	resp := suite.requestJsonWithAuth(api, nil, http.MethodPost, player1ID)

	// Then
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestStartGame_NotEnoughPlayers() {
	// Given - 只有 2 人
	hostID := suite.createAccount("Host8")
	player1ID := suite.createAccount("P8_1")

	gameID := suite.createGameRoom(hostID, 9)
	suite.joinGameRoom(gameID, player1ID)

	// When
	api := "/api/v1/games/" + strconv.Itoa(gameID) + "/start"
	resp := suite.requestJsonWithAuth(api, nil, http.MethodPost, hostID)

	// Then
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestStartGame_ExactlyThreePlayers() {
	// Given - 剛好 3 人
	hostID := suite.createAccount("Host9")
	player1ID := suite.createAccount("P9_1")
	player2ID := suite.createAccount("P9_2")

	gameID := suite.createGameRoom(hostID, 9)
	suite.joinGameRoom(gameID, player1ID)
	suite.joinGameRoom(gameID, player2ID)

	// When
	api := "/api/v1/games/" + strconv.Itoa(gameID) + "/start"
	resp := suite.requestJsonWithAuth(api, nil, http.MethodPost, hostID)

	// Then
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}
