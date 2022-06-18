package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lualfe/card-game/internal/usecase"

	"github.com/lualfe/card-game/internal/entity"

	"github.com/google/go-cmp/cmp"
)

type stubDeckManager struct {
	new       func(shuffle bool, cardCodes []string) entity.Deck
	open      func(id string) (entity.Deck, error)
	drawCards func(id string, amount int) ([]entity.Card, error)
}

func (s *stubDeckManager) DrawCards(id string, amount int) ([]entity.Card, error) {
	return s.drawCards(id, amount)
}

func (s *stubDeckManager) Open(id string) (entity.Deck, error) {
	return s.open(id)
}

func (s *stubDeckManager) New(shuffle bool, cardCodes []string) entity.Deck {
	return s.new(shuffle, cardCodes)
}

func Test_deckRoutes_newDeck(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       newDeckResponse
	}{
		{
			name:       "Success Not Shuffled",
			statusCode: http.StatusCreated,
			want: newDeckResponse{
				ID:        "id",
				Shuffled:  false,
				Remaining: 30,
			},
		},
		{
			name:       "Success Shuffled",
			statusCode: http.StatusCreated,
			want: newDeckResponse{
				ID:        "id",
				Shuffled:  true,
				Remaining: 30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodPost, "/deck/new", nil)

			d := &deckRoutes{
				deck: &stubDeckManager{
					new: func(shuffle bool, cardCodes []string) entity.Deck {
						return entity.Deck{
							ID:        tt.want.ID,
							Shuffled:  tt.want.Shuffled,
							Remaining: tt.want.Remaining,
							Cards:     []entity.Card{},
						}
					},
				},
			}
			d.newDeck(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			code := resp.StatusCode
			if code != tt.statusCode {
				t.Fatalf("deckRoutes.newDeck() | got status code %d, want %d", code, tt.statusCode)
			}

			var got newDeckResponse
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("deckRoutes.newDeck() | (-got +want):\n%s", diff)
			}
		})
	}
}

func Test_deckRoutes_openDeck(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       entity.Deck
		wantErr    error
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			want: entity.Deck{
				ID:        "id",
				Shuffled:  true,
				Remaining: 10,
				Cards:     entity.DefaultCards,
			},
		},
		{
			name:       "Unknown Error",
			statusCode: http.StatusInternalServerError,
			wantErr:    errors.New("random error"),
		},
		{
			name:       "Not Found Error",
			statusCode: http.StatusNotFound,
			wantErr:    usecase.DeckNotFoundErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodGet, "/v1/open/deck/id", nil)

			d := &deckRoutes{
				deck: &stubDeckManager{
					open: func(id string) (entity.Deck, error) {
						if tt.wantErr != nil {
							return entity.Deck{}, tt.wantErr
						}
						return tt.want, nil
					},
				},
			}
			d.openDeck(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			code := resp.StatusCode
			if code != tt.statusCode {
				t.Errorf("deckRoutes.openDeck() | got status code %d, want %d", code, tt.statusCode)
			}

			if tt.wantErr == nil {
				var got entity.Deck
				if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Fatalf("deckRoutes.openDeck() | (-got +want):\n%s", diff)
				}
			}
		})
	}
}

func Test_deckRoutes_drawCards(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    error
		want       drawCardsResp
	}{
		{
			name:       "Success",
			statusCode: http.StatusOK,
			want: drawCardsResp{
				Cards: []entity.Card{
					{
						Value: "ACE",
						Suit:  "SPADES",
						Code:  "AS",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r := httptest.NewRequest(http.MethodGet, "/v1/open/deck/id", nil)

			d := &deckRoutes{
				deck: &stubDeckManager{
					drawCards: func(id string, amount int) ([]entity.Card, error) {
						if tt.wantErr != nil {
							return nil, tt.wantErr
						}

						return tt.want.Cards, nil
					},
				},
			}
			d.drawCards(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			code := resp.StatusCode
			if code != tt.statusCode {
				t.Errorf("deckRoutes.openDeck() | got status code %d, want %d", code, tt.statusCode)
			}

			if tt.wantErr == nil {
				var got drawCardsResp
				if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Fatalf("deckRoutes.openDeck() | (-got +want):\n%s", diff)
				}
			}
		})
	}
}
