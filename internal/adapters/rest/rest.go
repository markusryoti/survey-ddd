package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/markusryoti/survey-ddd/internal/application/command"
	"github.com/markusryoti/survey-ddd/internal/application/query"
)

type SurveyHandler struct {
	CommandHandler *command.CommandHandler
	QueryHandler   *query.QueryHandler
}

func (h SurveyHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.index)
	r.Post("/surveys", h.CreateSurvey)
	r.Get("/surveys/{id}", h.GetSurvey)
	r.Post("/surveys/{id}/questions", h.AddQuestion)
}

func (h SurveyHandler) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = h.writeJson(w, map[string]string{
		"hello": "world",
	})
}

func (h SurveyHandler) writeJson(w http.ResponseWriter, body any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(body)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (h SurveyHandler) writeError(w http.ResponseWriter, status int, err ErrorResponse) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(err)
}
