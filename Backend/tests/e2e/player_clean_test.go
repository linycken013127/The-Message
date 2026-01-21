package e2e

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/usecase"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func (suite *CleanArchTestSuite) TestPlayCardE2E() {
	// Given
	api := "/api/v1/players/{player_id}/player-cards"
	game, _ := suite.gameUseCase.InitGame(context.TODO())

	// Fixed player count to 3 (InitIdentityCards only supports 3, 5-9 players)
	playerCount := 3

	// Fake players data
	var players []usecase.PlayerInfo
	for i := 0; i < playerCount; i++ {
		player := usecase.PlayerInfo{
			ID:   faker.UUIDDigit(),
			Name: faker.FirstName(),
		}
		players = append(players, player)
	}

	// Fake game data
	createGameRequest := usecase.CreateGameRequest{
		Players: players,
	}

	_ = suite.playerUseCase.InitPlayers(context.TODO(), game, createGameRequest)
	_ = suite.gameUseCase.InitDeck(context.TODO(), game)
	_ = suite.gameUseCase.DrawCardsForAllPlayers(context.TODO(), game)

	playerId := rand.Intn(playerCount) + 1

	// Get player's card
	cards, _ := suite.playerRepo.GetPlayerWithPlayerCards(context.TODO(), playerId)

	// Random card id
	num := rand.Intn(len(cards.PlayerCards))
	cardId := cards.PlayerCards[num].CardID

	// Set player to current player
	suite.gameUseCase.UpdateCurrentPlayer(context.TODO(), game, playerId)

	url := strings.ReplaceAll(api, "{player_id}", strconv.Itoa(playerId))
	req := PlayCardRequest{CardId: cardId}
	reqBody, _ := json.Marshal(req)

	res := suite.requestJson(url, reqBody, http.MethodPost)

	// Convert response body from json to map
	resBodyAsByteArray, _ := io.ReadAll(res.Body)
	resBody := make(map[string]interface{})
	_ = json.Unmarshal(resBodyAsByteArray, &resBody)

	// Then
	assert.Equal(suite.T(), 200, res.StatusCode)

	player, _ := suite.playerRepo.GetPlayerWithPlayerCards(context.TODO(), playerId)
	assert.Equal(suite.T(), len(cards.PlayerCards)-1, len(player.PlayerCards))
}

func (suite *CleanArchTestSuite) TestTransmitIntelligenceE2E() {
	api := "/api/v1/player/{player_id}/transmit-intelligence"
	game, _ := suite.gameUseCase.InitGame(context.TODO())

	// Fixed player count to 3 (InitIdentityCards only supports 3, 5-9 players)
	playerCount := 3

	// Fake players data
	var players []usecase.PlayerInfo
	for i := 0; i < playerCount; i++ {
		player := usecase.PlayerInfo{
			ID:   faker.UUIDDigit(),
			Name: faker.FirstName(),
		}
		players = append(players, player)
	}

	// Fake game data
	createGameRequest := usecase.CreateGameRequest{
		Players: players,
	}

	_ = suite.playerUseCase.InitPlayers(context.TODO(), game, createGameRequest)
	_ = suite.gameUseCase.InitDeck(context.TODO(), game)
	_ = suite.gameUseCase.DrawCardsForAllPlayers(context.TODO(), game)

	suite.T().Run("it can validate card id", func(t *testing.T) {
		playerId := rand.Intn(playerCount) + 1

		// Request only intelligence type
		url := strings.ReplaceAll(api, "{player_id}", strconv.Itoa(playerId))
		req := PlayCardRequest{}
		reqBody, _ := json.Marshal(req)

		res := suite.requestJson(url, reqBody, http.MethodPost)

		// Convert response body from json to map
		resBodyAsByteArray, _ := io.ReadAll(res.Body)
		resBody := make(map[string]interface{})
		_ = json.Unmarshal(resBodyAsByteArray, &resBody)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "Card not found", resBody["message"])
	})

	suite.T().Run("it can fail when player not found", func(t *testing.T) {
		playerId := math.MaxInt32
		cardId := rand.Intn(playerCount)

		url := strings.ReplaceAll(api, "{player_id}", strconv.Itoa(playerId))
		req := PlayCardRequest{CardId: cardId}
		reqBody, _ := json.Marshal(req)

		res := suite.requestJson(url, reqBody, http.MethodPost)

		// Convert response body from json to map
		resBodyAsByteArray, _ := io.ReadAll(res.Body)
		resBody := make(map[string]interface{})
		_ = json.Unmarshal(resBodyAsByteArray, &resBody)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "Player not found", resBody["message"])
	})

	suite.T().Run("it can fail when player card not found", func(t *testing.T) {
		playerId := rand.Intn(playerCount) + 1
		cardId := math.MaxInt32

		url := strings.ReplaceAll(api, "{player_id}", strconv.Itoa(playerId))
		req := PlayCardRequest{CardId: cardId}
		reqBody, _ := json.Marshal(req)

		res := suite.requestJson(url, reqBody, http.MethodPost)

		// Convert response body from json to map
		resBodyAsByteArray, _ := io.ReadAll(res.Body)
		resBody := make(map[string]interface{})
		_ = json.Unmarshal(resBodyAsByteArray, &resBody)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "Card not found", resBody["message"])
	})

	suite.T().Run("it can fail when game is end", func(t *testing.T) {
		playerId := rand.Intn(playerCount) + 1

		// Get player's card
		cards, _ := suite.playerRepo.GetPlayerWithPlayerCards(context.TODO(), playerId)

		// Random card id
		num := rand.Intn(len(cards.PlayerCards))
		cardId := cards.PlayerCards[num].CardID

		// Set player to current player
		suite.gameUseCase.UpdateCurrentPlayer(context.TODO(), game, playerId)

		// Set game status to end
		suite.gameUseCase.UpdateStatus(context.TODO(), game, entity.GameStatusEnd)

		url := strings.ReplaceAll(api, "{player_id}", strconv.Itoa(playerId))
		req := PlayCardRequest{CardId: cardId}
		reqBody, _ := json.Marshal(req)

		res := suite.requestJson(url, reqBody, http.MethodPost)

		// Convert response body from json to map
		resBodyAsByteArray, _ := io.ReadAll(res.Body)
		resBody := make(map[string]interface{})
		_ = json.Unmarshal(resBodyAsByteArray, &resBody)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "遊戲已結束", resBody["message"])

		// Recover game status to start
		suite.gameUseCase.UpdateStatus(context.TODO(), game, entity.GameStatusStart)
	})

	suite.T().Run("it can fail when not player's turn", func(t *testing.T) {
		playerId := rand.Intn(playerCount) + 1

		// Get player's card
		cards, _ := suite.playerRepo.GetPlayerWithPlayerCards(context.TODO(), playerId)

		// Random card id
		num := rand.Intn(len(cards.PlayerCards))
		cardId := cards.PlayerCards[num].CardID

		// Set other player to current player
		suite.gameUseCase.UpdateCurrentPlayer(context.TODO(), game, playerId-1)

		url := strings.ReplaceAll(api, "{player_id}", strconv.Itoa(playerId))
		req := PlayCardRequest{CardId: cardId}
		reqBody, _ := json.Marshal(req)

		res := suite.requestJson(url, reqBody, http.MethodPost)

		// Convert response body from json to map
		resBodyAsByteArray, _ := io.ReadAll(res.Body)
		resBody := make(map[string]interface{})
		_ = json.Unmarshal(resBodyAsByteArray, &resBody)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, "尚未輪到你出牌", resBody["message"])
	})

	suite.T().Run("it can success when valid card id", func(t *testing.T) {
		playerId := rand.Intn(playerCount) + 1

		// Get player's card
		cards, _ := suite.playerRepo.GetPlayerWithPlayerCards(context.TODO(), playerId)

		// Random card id
		num := rand.Intn(len(cards.PlayerCards))
		cardId := cards.PlayerCards[num].CardID

		// Set player to current player
		suite.gameUseCase.UpdateCurrentPlayer(context.TODO(), game, playerId)

		url := strings.ReplaceAll(api, "{player_id}", strconv.Itoa(playerId))
		req := PlayCardRequest{CardId: cardId}
		reqBody, _ := json.Marshal(req)

		res := suite.requestJson(url, reqBody, http.MethodPost)

		// Convert response body from json to map
		resBodyAsByteArray, _ := io.ReadAll(res.Body)
		resBody := make(map[string]interface{})
		_ = json.Unmarshal(resBodyAsByteArray, &resBody)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, true, resBody["result"])
	})
}
