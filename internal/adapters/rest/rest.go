package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/markusryoti/survey-ddd/internal/application/command"
	"github.com/markusryoti/survey-ddd/internal/application/query"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type SurveyHandler struct {
	CommandHandler *command.CommandHandler
	QueryHandler   *query.QueryHandler
}

func (h SurveyHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.index)
	r.Post("/surveys", h.CreateSurvey)
	r.Get("/surveys/{id}", h.GetSurvey)
}

func (h SurveyHandler) index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"hello": "world",
	})
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

	if req.TenantId == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	cmd := surveys.CreateSurveyCommand{
		Title:       req.Title,
		Description: req.Description,
		TenantId:    req.TenantId,
	}

	survey, err := h.CommandHandler.HandleCreateSurvey(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"survey_id": survey.Id.String(),
	})
}

type CreateSurveyRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	TenantId    string  `json:"tenantId"`
}

func (h SurveyHandler) GetSurvey(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	survey, err := h.QueryHandler.GetSurvey(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(survey)
}

type QuestionInput struct {
	Text         string                `json:"text"`
	QuestionType string                `json:"question_type"`
	Options      []QuestionOptionInput `json:"options"`
}

type QuestionOptionInput struct {
	Text string `json:"text"`
}
