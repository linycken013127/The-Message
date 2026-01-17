package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func (suite *CleanArchTestSuite) TestGetPlayerCards() {
	// given
	game := &entity.Game{}
	game, _ = suite.gameRepo.CreateGame(context.TODO(), game)
	
	player := entity.NewPlayer("player1", game.ID, "醬油", 1)
	player, _ = suite.playerRepo.CreatePlayer(context.TODO(), player)

	playerCard := entity.NewPlayerCard(player.ID, game.ID, 1, entity.PlayerCardTypeHand)
	_, err := suite.playerCardRepo.CreatePlayerCard(context.TODO(), playerCard)
	if err != nil {
		panic(err)
	}

	// when
	api := "/api/v1/player/1/player-cards/"
	resp := suite.requestJson(api, nil, http.MethodGet)
	response := suite.responseJson(resp)
	
	// then
	assert.Equal(suite.T(), 200, resp.StatusCode)

	playerCards, ok := response["player_cards"]
	if !ok {
		fmt.Println("Error: player_cards is not of type []interface{}")
		return
	}

	if str, ok := playerCards.(string); ok {
		var slice []map[string]interface{}
		if err := json.Unmarshal([]byte(str), &slice); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// Range over the slice
		for _, item := range slice {
			fmt.Println(item["id"], item["name"], item["color"])
			for key, value := range item {
				if value == nil {
					suite.T().Errorf("Field %s is nil", key)
				}
			}
		}
	}
}
