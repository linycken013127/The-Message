package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/stretchr/testify/assert"
)

// RegisterPlayerRequest 註冊玩家請求
type RegisterPlayerRequest struct {
	PlayerName string `json:"playerName"`
}

func (suite *CleanArchTestSuite) TestRegisterPlayer_Success() {
	// Given
	api := "/api/v1/players"
	req := RegisterPlayerRequest{PlayerName: "Alice"}
	jsonBody, _ := json.Marshal(req)

	// When
	resp := suite.requestJson(api, jsonBody, http.MethodPost)

	// Then
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	assert.NotNil(suite.T(), responseJson["playerId"], "回應應該包含 playerId")
	assert.Equal(suite.T(), "Alice", responseJson["playerName"], "回應應該包含正確的 playerName")

	// 驗證資料庫中存在該玩家
	playerId := int(responseJson["playerId"].(float64))
	account, err := suite.accountRepo.GetAccountById(context.TODO(), playerId)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Alice", account.Name)
}

func (suite *CleanArchTestSuite) TestRegisterPlayer_DuplicateName() {
	// Given - 先建立一個玩家
	api := "/api/v1/players"
	req1 := RegisterPlayerRequest{PlayerName: "Alice"}
	jsonBody1, _ := json.Marshal(req1)
	suite.requestJson(api, jsonBody1, http.MethodPost)

	// When - 嘗試使用相同名稱建立另一個玩家
	req2 := RegisterPlayerRequest{PlayerName: "Alice"}
	jsonBody2, _ := json.Marshal(req2)
	resp := suite.requestJson(api, jsonBody2, http.MethodPost)

	// Then - 應該失敗
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	assert.NotNil(suite.T(), responseJson["message"])
}

func (suite *CleanArchTestSuite) TestRegisterPlayer_NameTooLong() {
	// Given - 名稱超過 10 個字元
	api := "/api/v1/players"
	req := RegisterPlayerRequest{PlayerName: "12345678901"} // 11 字元
	jsonBody, _ := json.Marshal(req)

	// When
	resp := suite.requestJson(api, jsonBody, http.MethodPost)

	// Then - 應該失敗
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	assert.NotNil(suite.T(), responseJson["message"])
}

func (suite *CleanArchTestSuite) TestRegisterPlayer_NameExactly10Chars() {
	// Given - 名稱剛好 10 個字元
	api := "/api/v1/players"
	req := RegisterPlayerRequest{PlayerName: "1234567890"} // 10 字元
	jsonBody, _ := json.Marshal(req)

	// When
	resp := suite.requestJson(api, jsonBody, http.MethodPost)

	// Then - 應該成功
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestRegisterPlayer_EmptyName() {
	// Given - 空名稱
	api := "/api/v1/players"
	req := RegisterPlayerRequest{PlayerName: ""}
	jsonBody, _ := json.Marshal(req)

	// When
	resp := suite.requestJson(api, jsonBody, http.MethodPost)

	// Then - 應該失敗
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestRegisterPlayer_NameWithSpaces() {
	// Given - 名稱包含空格
	api := "/api/v1/players"
	req := RegisterPlayerRequest{PlayerName: "Al ice"}
	jsonBody, _ := json.Marshal(req)

	// When
	resp := suite.requestJson(api, jsonBody, http.MethodPost)

	// Then - 應該成功（空格是合法字元）
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestRegisterPlayer_DuplicateNameCaseInsensitive() {
	// Given - 先建立一個玩家 "alice"
	api := "/api/v1/players"
	req1 := RegisterPlayerRequest{PlayerName: "alice"}
	jsonBody1, _ := json.Marshal(req1)
	suite.requestJson(api, jsonBody1, http.MethodPost)

	// When - 嘗試使用 "Alice" (不同大小寫) 建立另一個玩家
	req2 := RegisterPlayerRequest{PlayerName: "Alice"}
	jsonBody2, _ := json.Marshal(req2)
	resp := suite.requestJson(api, jsonBody2, http.MethodPost)

	// Then - 應該失敗（名稱比較不區分大小寫）
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *CleanArchTestSuite) TestRegisterPlayer_TrimWhitespace() {
	// Given - 名稱前後有空白
	api := "/api/v1/players"
	req := RegisterPlayerRequest{PlayerName: "  Bob  "}
	jsonBody, _ := json.Marshal(req)

	// When
	resp := suite.requestJson(api, jsonBody, http.MethodPost)

	// Then - 應該成功，並且名稱被 trim
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	responseJson := suite.responseJson(resp)
	// 驗證名稱已被 trim
	playerName := responseJson["playerName"].(string)
	assert.Equal(suite.T(), "Bob", strings.TrimSpace(playerName))
}
