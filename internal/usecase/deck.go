package usecase

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/lualfe/card-game/internal/usecase/repo"

	"github.com/google/uuid"

	"github.com/lualfe/card-game/internal/entity"
)

var (
	// DeckNotFoundErr happens when a deck can't be found in the repo.
	DeckNotFoundErr = errors.New("deck not found")
)

// Deck is a use case to manage the game deck.
type Deck struct {
	deckRepo DeckRepo
	shuffler func([]entity.Card)
}

// NewDeckManager creates a new Deck.
func NewDeckManager(store DeckRepo) *Deck {
	return &Deck{
		deckRepo: store,
		shuffler: func(cards []entity.Card) {
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(cards), func(i, j int) {
				cards[i], cards[j] = cards[j], cards[i]
			})
		},
	}
}

// New generates a new entity.Deck.
func (d *Deck) New(shuffle bool, cardCodes []string) entity.Deck {
	deckCards := entity.DefaultCards

	if cardCodes != nil && len(cardCodes) > 0 {
		deckCards = []entity.Card{}
		for _, cc := range cardCodes {
			for _, dc := range entity.DefaultCards {
				if cc == dc.Code {
					deckCards = append(deckCards, dc)
				}
			}
		}
	}

	if shuffle {
		d.shuffler(deckCards)
	}

	deck := entity.Deck{
		ID:        uuid.New().String(),
		Shuffled:  shuffle,
		Remaining: len(deckCards),
		Cards:     deckCards,
	}

	d.deckRepo.Save(deck)

	return deck
}

// Open returns a deck or an error in case the
// deck can't be found.
func (d *Deck) Open(id string) (entity.Deck, error) {
	deck, err := d.deckRepo.Get(id)
	if err != nil {
		if errors.Is(err, repo.DeckNotFoundErr) {
			return entity.Deck{}, fmt.Errorf("%w with id %s", DeckNotFoundErr, id)
		}
		return entity.Deck{}, err
	}

	return deck, nil
}

// DrawCards gets cards from the top of the deck.
func (d *Deck) DrawCards(id string, amount int) ([]entity.Card, error) {
	deck, err := d.deckRepo.Get(id)
	if err != nil {
		if errors.Is(err, repo.DeckNotFoundErr) {
			return nil, fmt.Errorf("%w with id %s", DeckNotFoundErr, id)
		}
		return nil, err
	}

	if amount > deck.Remaining {
		amount = deck.Remaining
	}

	if amount < 0 {
		amount = 1
	}

	cards := deck.Cards[:amount]
	deck.Cards = deck.Cards[amount:]
	deck.Remaining = len(deck.Cards)

	d.deckRepo.Save(deck)

	return cards, nil
}
