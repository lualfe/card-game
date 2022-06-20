package v1

import (
	"github.com/go-chi/chi/v5"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/lualfe/card-game/docs"
	"github.com/lualfe/card-game/internal/usecase"
)

// @title        Decks API
// @version      1.0
// @description  This is a server to handle a cards deck.

// @host      localhost:8080
// @BasePath  /v1

// StartRoutes starts the application routes.
func StartRoutes(m *chi.Mux, deck usecase.DeckManager) {
	m.Mount("/swagger", httpSwagger.WrapHandler)
	createDeckRoutes(m, deck)
}
