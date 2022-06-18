package repo

import (
	"errors"
	"fmt"

	"github.com/lualfe/card-game/internal/entity"
)

// DeckNotFoundErr happens when a deck is not found
// in the repo.
var DeckNotFoundErr = errors.New("deck not found")

// Deck repo.
type Deck map[string]entity.Deck

// Save saves a deck to the store.
func (d Deck) Save(deck entity.Deck) {
	d[deck.ID] = deck
}

// Get retrieves a deck from its ID.
func (d Deck) Get(id string) (entity.Deck, error) {
	deck, ok := d[id]
	if !ok {
		return entity.Deck{}, fmt.Errorf("%w with ID %s", DeckNotFoundErr, id)
	}
	return deck, nil
}
