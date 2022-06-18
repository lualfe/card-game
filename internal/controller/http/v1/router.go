package v1

import (
	"github.com/go-chi/chi/v5"

	"github.com/lualfe/card-game/internal/usecase"
)

func Routes(m *chi.Mux, deck usecase.DeckManager) {
	createDeckRoutes(m, deck)
}
