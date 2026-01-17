package dto

import (
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
)

// GameDTO 遊戲資料傳輸物件
type GameDTO struct {
	ID              int          `json:"id"`
	Token           string       `json:"token"`
	Status          string       `json:"status"`
	CurrentPlayerID int          `json:"currentPlayerId"`
	Players         []PlayerDTO  `json:"players"`
	CreatedAt       time.Time    `json:"createdAt"`
}

// ToGameDTO 將 Entity 轉換為 DTO
func ToGameDTO(game *entity.Game) *GameDTO {
	if game == nil {
		return nil
	}

	dto := &GameDTO{
		ID:              game.ID,
		Token:           game.Token,
		Status:          game.Status,
		CurrentPlayerID: game.CurrentPlayerID,
		CreatedAt:       game.CreatedAt,
		Players:         make([]PlayerDTO, 0, len(game.Players)),
	}

	for _, p := range game.Players {
		dto.Players = append(dto.Players, *ToPlayerDTO(&p))
	}

	return dto
}

// PlayerDTO 玩家資料傳輸物件
type PlayerDTO struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	GameID       int           `json:"gameId"`
	IdentityCard string        `json:"identityCard,omitempty"`
	Status       string        `json:"status"`
	OrderNumber  int           `json:"orderNumber"`
	PlayerCards  []PlayerCardDTO `json:"playerCards,omitempty"`
}

// ToPlayerDTO 將 Entity 轉換為 DTO
func ToPlayerDTO(player *entity.Player) *PlayerDTO {
	if player == nil {
		return nil
	}

	dto := &PlayerDTO{
		ID:           player.ID,
		Name:         player.Name,
		GameID:       player.GameID,
		IdentityCard: player.IdentityCard,
		Status:       player.Status,
		OrderNumber:  player.OrderNumber,
		PlayerCards:  make([]PlayerCardDTO, 0, len(player.PlayerCards)),
	}

	for _, c := range player.PlayerCards {
		dto.PlayerCards = append(dto.PlayerCards, *ToPlayerCardDTO(&c))
	}

	return dto
}

// CardDTO 卡片資料傳輸物件
type CardDTO struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Color            string `json:"color"`
	IntelligenceType int    `json:"intelligenceType"`
}

// ToCardDTO 將 Entity 轉換為 DTO
func ToCardDTO(card *entity.Card) *CardDTO {
	if card == nil {
		return nil
	}

	return &CardDTO{
		ID:               card.ID,
		Name:             card.Name,
		Color:            card.Color,
		IntelligenceType: card.IntelligenceType,
	}
}

// PlayerCardDTO 玩家手牌資料傳輸物件
type PlayerCardDTO struct {
	ID       int     `json:"id"`
	PlayerID int     `json:"playerId"`
	GameID   int     `json:"gameId"`
	CardID   int     `json:"cardId"`
	Type     string  `json:"type"`
	Card     *CardDTO `json:"card,omitempty"`
}

// ToPlayerCardDTO 將 Entity 轉換為 DTO
func ToPlayerCardDTO(playerCard *entity.PlayerCard) *PlayerCardDTO {
	if playerCard == nil {
		return nil
	}

	dto := &PlayerCardDTO{
		ID:       playerCard.ID,
		PlayerID: playerCard.PlayerID,
		GameID:   playerCard.GameID,
		CardID:   playerCard.CardID,
		Type:     playerCard.Type,
	}

	if playerCard.Card.ID != 0 {
		dto.Card = ToCardDTO(&playerCard.Card)
	}

	return dto
}

// DeckDTO 牌組資料傳輸物件
type DeckDTO struct {
	ID     int `json:"id"`
	GameID int `json:"gameId"`
	CardID int `json:"cardId"`
}

// ToDeckDTO 將 Entity 轉換為 DTO
func ToDeckDTO(deck *entity.Deck) *DeckDTO {
	if deck == nil {
		return nil
	}

	return &DeckDTO{
		ID:     deck.ID,
		GameID: deck.GameID,
		CardID: deck.CardID,
	}
}
