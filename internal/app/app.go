package app

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/lualfe/card-game/internal/usecase"
	"github.com/lualfe/card-game/internal/usecase/repo"

	v1 "github.com/lualfe/card-game/internal/controller/http/v1"
)

// Run create all the main objects and run the
// application.
func Run() {
	m := chi.NewRouter()

	deckRepo := make(repo.Deck)
	dm := usecase.NewDeckManager(deckRepo)

	v1.Routes(m, dm)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", m))
}
