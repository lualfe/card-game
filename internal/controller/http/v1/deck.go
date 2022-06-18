package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/lualfe/card-game/internal/entity"

	"github.com/go-chi/chi/v5"

	"github.com/lualfe/card-game/internal/controller/http/response"
	"github.com/lualfe/card-game/internal/usecase"
)

func createDeckRoutes(m *chi.Mux, deck usecase.DeckManager) {
	dr := &deckRoutes{deck}

	m.Route("/v1/decks", func(r chi.Router) {
		r.Post("/", dr.newDeck)
		r.Get("/{deckID}", dr.openDeck)
		r.Get("/withdrawals/{deckID}", dr.drawCards)
	})
}

type deckRoutes struct {
	deck usecase.DeckManager
}

type newDeckResponse struct {
	ID        string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

func (d *deckRoutes) newDeck(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	shuffle := false
	if q.Get("shuffle") == "true" {
		shuffle = true
	}

	var cardCodes []string
	if cards := q.Get("cards"); cards != "" {
		cardCodes = strings.Split(cards, ",")
	}

	deck := d.deck.New(shuffle, cardCodes)

	resp := newDeckResponse{
		ID:        deck.ID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining,
	}

	response.JSON(w, resp, http.StatusCreated)
}

func (d *deckRoutes) openDeck(w http.ResponseWriter, r *http.Request) {
	deckID := chi.URLParam(r, "deckID")

	deck, err := d.deck.Open(deckID)
	if err != nil {
		if errors.Is(err, usecase.DeckNotFoundErr) {
			response.JSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		response.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.JSON(w, deck, http.StatusOK)
}

type drawCardsResp struct {
	Cards []entity.Card `json:"cards"`
}

func (d *deckRoutes) drawCards(w http.ResponseWriter, r *http.Request) {
	deckID := chi.URLParam(r, "deckID")
	amount := 1
	if am := r.URL.Query().Get("amount"); am != "" {
		if v, err := strconv.Atoi(am); err == nil {
			amount = v
		}
	}

	cards, err := d.deck.DrawCards(deckID, amount)
	if err != nil {
		response.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := drawCardsResp{
		Cards: cards,
	}

	response.JSON(w, resp, http.StatusOK)
}
