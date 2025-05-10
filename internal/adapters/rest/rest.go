package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/markusryoti/survey-ddd/internal/application/command"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type SurveyHandler struct {
	CommandHandler *command.SurveyCommandHandler
}

func (h SurveyHandler) RegisterRoutes(r chi.Router) {
	r.Post("/surveys", h.CreateSurvey)
}

func (h SurveyHandler) CreateSurvey(w http.ResponseWriter, r *http.Request) {
	var req CreateSurveyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	cmd := surveys.CreateSurveyCommand{
		Title:       req.Title,
		Description: req.Description,
	}

	survey, err := h.CommandHandler.HandleCreateSurvey(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"survey_id": string(survey.Id.String()),
	})
}

type CreateSurveyRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type QuestionInput struct {
	Text         string                `json:"text"`
	QuestionType string                `json:"question_type"`
	Options      []QuestionOptionInput `json:"options"`
}

type QuestionOptionInput struct {
	Text string `json:"text"`
}
