package usecase

import (
	"context"
	"math/rand"
	"time"

	"github.com/Game-as-a-Service/The-Message/internal/domain/entity"
	"github.com/Game-as-a-Service/The-Message/internal/domain/repository"
)

// DeckUseCase 牌組用例介面
type DeckUseCase interface {
	// InitDeck 初始化牌組
	InitDeck(ctx context.Context, game *entity.Game) error
	// ShuffleDeck 洗牌
	ShuffleDeck(cards []*entity.Card) []*entity.Card
	// CreateDeck 建立牌組
	CreateDeck(ctx context.Context, deck *entity.Deck) (*entity.Deck, error)
	// GetDecksByGameId 取得牌組
	GetDecksByGameId(ctx context.Context, gameID int) ([]*entity.Deck, error)
	// DeleteDeckFromGame 從遊戲中刪除牌組
	DeleteDeckFromGame(ctx context.Context, id int) error
}

// deckUseCase 牌組用例實作
type deckUseCase struct {
	cardUseCase CardUseCase
	deckRepo    repository.DeckRepository
}

// DeckUseCaseOptions 牌組用例選項
type DeckUseCaseOptions struct {
	CardUseCase CardUseCase
	DeckRepo    repository.DeckRepository
}

// NewDeckUseCase 建立牌組用例
func NewDeckUseCase(opts *DeckUseCaseOptions) DeckUseCase {
	return &deckUseCase{
		cardUseCase: opts.CardUseCase,
		deckRepo:    opts.DeckRepo,
	}
}

// InitDeck 初始化牌組
func (uc *deckUseCase) InitDeck(ctx context.Context, game *entity.Game) error {
	cards, err := uc.cardUseCase.GetCards(ctx)
	if err != nil {
		return err
	}

	cards = uc.ShuffleDeck(cards)

	for _, card := range cards {
		_, err := uc.deckRepo.CreateDeck(ctx, entity.NewDeck(game.ID, card.ID))
		if err != nil {
			return err
		}
	}

	return nil
}

// ShuffleDeck 洗牌
func (uc *deckUseCase) ShuffleDeck(cards []*entity.Card) []*entity.Card {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return cards
}

// CreateDeck 建立牌組
func (uc *deckUseCase) CreateDeck(ctx context.Context, deck *entity.Deck) (*entity.Deck, error) {
	deck, err := uc.deckRepo.CreateDeck(ctx, deck)
	if err != nil {
		return nil, err
	}
	return deck, nil
}

// GetDecksByGameId 取得牌組
func (uc *deckUseCase) GetDecksByGameId(ctx context.Context, gameID int) ([]*entity.Deck, error) {
	decks, err := uc.deckRepo.GetDecksByGameId(ctx, gameID)
	if err != nil {
		return nil, err
	}
	return decks, nil
}

// DeleteDeckFromGame 從遊戲中刪除牌組
func (uc *deckUseCase) DeleteDeckFromGame(ctx context.Context, id int) error {
	return uc.deckRepo.DeleteDeck(ctx, id)
}
