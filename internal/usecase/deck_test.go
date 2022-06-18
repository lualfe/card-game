package usecase

import (
	"errors"
	"testing"

	"github.com/lualfe/card-game/internal/usecase/repo"

	"github.com/google/go-cmp/cmp"

	"github.com/lualfe/card-game/internal/entity"
)

type stubDeckStore struct {
	get func(id string) (entity.Deck, error)
}

func (s *stubDeckStore) Get(id string) (entity.Deck, error) {
	return s.get(id)
}

func (s *stubDeckStore) Save(deck entity.Deck) {}

func TestDeck_New(t *testing.T) {
	customDeck := []entity.Card{
		{
			Value: "ACE",
			Suit:  "SPADES",
			Code:  "AS",
		},
		{
			Value: "2",
			Suit:  "SPADES",
			Code:  "2S",
		},
	}

	nonexistentCards := []entity.Card{
		{
			Value: "ACE",
			Suit:  "SPADES",
			Code:  "AS",
		},
		{
			Value: "2",
			Suit:  "SPADES",
			Code:  "2S",
		},
		{
			Value: "bad",
			Suit:  "SPADES",
			Code:  "bad code",
		},
	}

	tests := []struct {
		name      string
		cardCodes []string
		want      entity.Deck
	}{
		{
			name: "Shuffled Default Cards",
			want: entity.Deck{
				Shuffled:  true,
				Remaining: 52,
				Cards:     entity.DefaultCards,
			},
		},
		{
			name: "Not Shuffled Default Cards",
			want: entity.Deck{
				Shuffled:  false,
				Remaining: 52,
				Cards:     entity.DefaultCards,
			},
		},
		{
			name: "Custom Cards",
			cardCodes: func() []string {
				var codes []string
				for _, c := range customDeck {
					codes = append(codes, c.Code)
				}
				return codes
			}(),
			want: entity.Deck{
				Shuffled:  true,
				Remaining: len(customDeck),
				Cards:     customDeck,
			},
		},
		{
			name: "Ignore Nonexistent Cards",
			cardCodes: func() []string {
				var codes []string
				for _, c := range nonexistentCards {
					codes = append(codes, c.Code)
				}
				return codes
			}(),
			want: entity.Deck{
				Shuffled:  true,
				Remaining: 2,
				Cards:     nonexistentCards[:len(nonexistentCards)-1],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shufflerCalled := false
			d := &Deck{
				deckStore: &stubDeckStore{},
				shuffler: func(cards []entity.Card) {
					shufflerCalled = true
				},
			}

			got := d.New(tt.want.Shuffled, tt.cardCodes)
			if got.ID == "" {
				t.Error("Deck.New() | got empty ID")
			}
			if tt.want.Shuffled && !shufflerCalled {
				t.Error("Deck.New() | wants shuffled cards")
			}
			tt.want.ID = got.ID
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("Deck.New() | (-got +want):\n%s", diff)
			}
		})
	}
}

func TestDeck_Open(t *testing.T) {
	unknownErr := errors.New("error")

	tests := []struct {
		name    string
		want    entity.Deck
		getErr  error
		wantErr error
	}{
		{
			name: "Success",
			want: entity.Deck{
				ID:        "id",
				Shuffled:  true,
				Remaining: 52,
				Cards:     entity.DefaultCards,
			},
		},
		{
			name:    "Unknown Error",
			getErr:  unknownErr,
			wantErr: unknownErr,
		},
		{
			name:    "Not Found Error",
			getErr:  repo.DeckNotFoundErr,
			wantErr: DeckNotFoundErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deck{
				deckStore: &stubDeckStore{
					get: func(id string) (entity.Deck, error) {
						if tt.wantErr != nil {
							return entity.Deck{}, tt.getErr
						}
						return tt.want, nil
					},
				},
			}

			got, err := d.Open(tt.want.ID)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Deck.Open() | got error %v, want nil", err)
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Fatalf("Deck.Open() | (-got +want):\n%s", diff)
				}
			} else {
				if err == nil {
					t.Errorf("Deck.Open() | got error nil, want not nil")
				}

				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Deck.Open() | got error %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestDeck_DrawCards(t *testing.T) {
	unknownErr := errors.New("error")

	tests := []struct {
		name    string
		deckID  string
		amount  int
		want    []entity.Card
		getErr  error
		wantErr error
	}{
		{
			name:   "Success",
			deckID: "id",
			amount: 1,
			want: []entity.Card{
				{
					Value: "ACE",
					Suit:  "SPADES",
					Code:  "AS",
				},
			},
		},
		{
			name:    "Deck Not Found",
			deckID:  "id",
			amount:  1,
			getErr:  repo.DeckNotFoundErr,
			wantErr: DeckNotFoundErr,
		},
		{
			name:    "Unknown Error",
			deckID:  "id",
			amount:  1,
			getErr:  unknownErr,
			wantErr: unknownErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Deck{
				deckStore: &stubDeckStore{
					get: func(id string) (entity.Deck, error) {
						if tt.wantErr != nil {
							return entity.Deck{}, tt.getErr
						}
						return entity.Deck{
							ID:        "id",
							Shuffled:  false,
							Remaining: len(tt.want),
							Cards:     tt.want,
						}, nil
					},
				},
			}
			got, err := d.DrawCards(tt.deckID, tt.amount)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Deck.DrawCards() | error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("Deck.DrawCards() | (-got +want):\n%s", diff)
				}
			} else {
				if err == nil {
					t.Error("Deck.DrawCards() | got error nil, want not nil")
				}

				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Deck.DrawCards() | got error %v, want %v", err, tt.wantErr)
				}
			}
		})
	}
}
