package usecase

import "github.com/lualfe/card-game/internal/entity"

// DeckManager is the interface for deck operations.
type DeckManager interface {
	New(shuffle bool, cardCodes []string) entity.Deck
	Open(id string) (entity.Deck, error)
	DrawCards(id string, amount int) ([]entity.Card, error)
}

// DeckStore is the interface for the deck store.
type DeckStore interface {
	Save(deck entity.Deck)
	Get(id string) (entity.Deck, error)
}
