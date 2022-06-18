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

// newDeck godoc
// @Summary      Creates a new deck.
// @Description  Creates a new deck with cards.
// @Produce      json
// @Param        shuffle  query     bool    false  "Activate or deactivate cards shuffling."                                                                      default(false)
// @Param        cards    query     string  false  "Comma separated card codes to create a custom deck. If not sent, the regular 52 cards deck will be created."  example(AS,2S)
// @Success      200      {object}  newDeckResponse
// @Router       /decks [post]
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

// openDeck godoc
// @Summary      Opens a deck.
// @Description  Opens a deck, showing all its cards.
// @Produce      json
// @Param        id   path      string  true  "Deck id"
// @Success      200  {object}  entity.Deck
// @Failure      404     {object}  response.Error
// @Failure      500     {object}  response.Error
// @Router       /decks/{id} [get]
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

// drawCards godoc
// @Summary      Draw cards from a deck.
// @Description  Draw an amount of cards given a deck.
// @Produce      json
// @Param        id      path      string  true   "Deck id"
// @Param        amount  query     int     false  "Amount of cards to draw"  default(1)
// @Success      200     {object}  drawCardsResp
// @Failure      404  {object}  response.Error
// @Failure      500  {object}  response.Error
// @Router       /decks/withdrawals/{id} [get]
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
		if errors.Is(err, usecase.DeckNotFoundErr) {
			response.JSONError(w, err.Error(), http.StatusNotFound)
			return
		}

		response.JSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := drawCardsResp{
		Cards: cards,
	}

	response.JSON(w, resp, http.StatusOK)
}
