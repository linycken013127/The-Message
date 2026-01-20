package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/stretchr/testify/assert"
)

// ============ 認證機制測試 ============

func (suite *CleanArchTestSuite) TestAuth_ValidAccountID() {
	// Given - 建立有效帳號
	accountID := suite.createAccount("AuthTest1")

	// When - 使用有效帳號建立遊戲房
	api := "/api/v1/games"
	req := CreateGameRoomRequest{MaxPlayers: 9}
	jsonBody, _ := json.Marshal(req)
	resp := suite.requestJsonWithAuth(api, jsonBody, http.MethodPost, accountID)

	// Then - 應該成功
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestAuth_MissingHeader() {
	// When - 不帶 X-Account-ID header
	api := "/api/v1/games"
	req := CreateGameRoomRequest{MaxPlayers: 9}
	jsonBody, _ := json.Marshal(req)
	resp := suite.requestJson(api, jsonBody, http.MethodPost)

	// Then - 應該回傳 401
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	assert.Contains(suite.T(), responseJson["message"], "未提供認證資訊")
}

func (suite *CleanArchTestSuite) TestAuth_InvalidFormat() {
	// When - X-Account-ID 格式不正確（非數字）
	api := "/api/v1/games"
	req := CreateGameRoomRequest{MaxPlayers: 9}
	jsonBody, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest(http.MethodPost, suite.server.URL+api, bytes.NewBuffer(jsonBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Account-ID", "invalid")

	client := &http.Client{}
	resp, _ := client.Do(httpReq)

	// Then - 應該回傳 401
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	assert.Contains(suite.T(), responseJson["message"], "無效的認證資訊")
}

func (suite *CleanArchTestSuite) TestAuth_NonExistentAccount() {
	// When - 使用不存在的帳號 ID
	api := "/api/v1/games"
	req := CreateGameRoomRequest{MaxPlayers: 9}
	jsonBody, _ := json.Marshal(req)

	// 使用一個很大的 ID，確保不存在
	resp := suite.requestJsonWithAuth(api, jsonBody, http.MethodPost, 999999)

	// Then - 應該回傳 401
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	assert.Contains(suite.T(), responseJson["message"], "帳號不存在")
}

func (suite *CleanArchTestSuite) TestAuth_NegativeAccountID() {
	// When - 使用負數帳號 ID
	api := "/api/v1/games"
	req := CreateGameRoomRequest{MaxPlayers: 9}
	jsonBody, _ := json.Marshal(req)
	resp := suite.requestJsonWithAuth(api, jsonBody, http.MethodPost, -1)

	// Then - 應該回傳 401（帳號不存在）
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestAuth_ZeroAccountID() {
	// When - 使用 0 作為帳號 ID
	api := "/api/v1/games"
	req := CreateGameRoomRequest{MaxPlayers: 9}
	jsonBody, _ := json.Marshal(req)
	resp := suite.requestJsonWithAuth(api, jsonBody, http.MethodPost, 0)

	// Then - 應該回傳 401（帳號不存在）
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestAuth_JoinGameRoom_RequiresAuth() {
	// Given - 建立一個遊戲房
	hostID := suite.createAccount("AuthHost")
	gameID := suite.createGameRoom(hostID, 9)

	// When - 不帶認證嘗試加入
	api := "/api/v1/games/" + strconv.Itoa(gameID) + "/join"
	resp := suite.requestJson(api, nil, http.MethodPost)

	// Then - 應該回傳 401
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestAuth_StartGame_RequiresAuth() {
	// Given - 建立一個遊戲房
	hostID := suite.createAccount("AuthHost2")
	gameID := suite.createGameRoom(hostID, 9)

	// When - 不帶認證嘗試開始遊戲
	api := "/api/v1/games/" + strconv.Itoa(gameID) + "/start"
	resp := suite.requestJson(api, nil, http.MethodPost)

	// Then - 應該回傳 401
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}
