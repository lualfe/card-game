package usecase

import "github.com/lualfe/card-game/internal/entity"

// DeckManager is the interface for deck operations.
type DeckManager interface {
	New(shuffle bool, cardCodes []string) entity.Deck
	Open(id string) (entity.Deck, error)
	DrawCards(id string, amount int) ([]entity.Card, error)
}

// DeckRepo is the interface for the deck store.
type DeckRepo interface {
	Save(deck entity.Deck)
	Get(id string) (entity.Deck, error)
}
