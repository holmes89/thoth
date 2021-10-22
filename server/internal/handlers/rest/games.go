package rest

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/holmes89/thoth/internal"
	"github.com/rs/zerolog/log"
)

type gameHandler struct {
	svc GameService
}

func MakeGameHandler(mr *mux.Router, svc GameService) {
	r := mr.PathPrefix("/games").Subrouter()

	h := gameHandler{
		svc: svc,
	}

	r.HandleFunc("", h.FindAll).Methods("GET", "OPTIONS")
	r.HandleFunc("/", h.FindAll).Methods("GET", "OPTIONS")
	r.HandleFunc("/{id}", h.FindByID).Methods("GET", "OPTIONS")
}

func (h *gameHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resp, err := h.svc.ListGames(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to find games")
		http.Error(w, "unable to find games", http.StatusInternalServerError)
		return
	}
	EncodeResponse(w, resp, err)
}

func (h *gameHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ctx := r.Context()
	id := vars["id"]

	resp, err := h.svc.FindGame(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("unable to find game")
		http.Error(w, "unable to find game", http.StatusInternalServerError)
		return
	}
	EncodeResponse(w, resp, err)
}

type GameService interface {
	ListGames(ctx context.Context) ([]internal.Game, error)
	FindGame(context.Context, string) (internal.Game, error)
}
