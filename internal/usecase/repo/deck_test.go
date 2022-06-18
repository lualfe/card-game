package repo

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/lualfe/card-game/internal/entity"
)

func TestDeck_Save(t *testing.T) {
	tests := []struct {
		name string
		ent  entity.Deck
	}{
		{
			name: "Success",
			ent: entity.Deck{
				ID:        "id",
				Shuffled:  true,
				Remaining: 52,
				Cards:     entity.DefaultCards,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deckStore := make(Deck)
			deckStore.Save(tt.ent)

			got, ok := deckStore[tt.ent.ID]
			if !ok {
				t.Fatalf("Deck.Save() | saved deck not found in the store")
			}

			if diff := cmp.Diff(got, tt.ent); diff != "" {
				t.Fatalf("Deck.Save() | (-got +want):\n%s", diff)
			}
		})
	}
}

func TestDeck_Get(t *testing.T) {
	tests := []struct {
		name string
		want entity.Deck
	}{
		{
			name: "Success",
			want: entity.Deck{
				ID:        "id",
				Shuffled:  false,
				Remaining: 52,
				Cards:     entity.DefaultCards,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deckStore := make(Deck)
			deckStore[tt.want.ID] = tt.want
			got, err := deckStore.Get(tt.want.ID)
			if err != nil {
				t.Errorf("Deck.Get() | got error %v, want nil", err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("Deck.Get() | (-got +want):\n%s", diff)
			}
		})
	}
}

func TestDeck_Get_Error(t *testing.T) {
	deckStore := make(Deck)
	_, err := deckStore.Get("id")
	if err == nil {
		t.Errorf("Deck.Get() | got error %v, want nil", err)
	}
}
