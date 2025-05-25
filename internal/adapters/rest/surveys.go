package rest

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

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

	survey, err := h.CommandHandler.CreateSurvey(r.Context(), cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = h.writeJson(w, map[string]string{
		"surveyId": survey.Id.String(),
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

	_ = h.writeJson(w, survey)
}

type AddQuestionRequest struct {
	Title           string   `json:"title"`
	Description     *string  `json:"description"`
	AllowMultiple   bool     `json:"allowMultiple"`
	QuestionOptions []string `json:"questionOptions"`
}

func (h SurveyHandler) AddQuestion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req AddQuestionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.CommandHandler.AddQuestion(r.Context(), surveys.AddQuestionCommand{
		SurveyId:        id,
		Title:           req.Title,
		Description:     req.Description,
		AllowMultiple:   req.AllowMultiple,
		QuestionOptions: req.QuestionOptions,
	})

	if err != nil {
		h.writeError(w, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type QuestionInput struct {
	Text         string                `json:"text"`
	QuestionType string                `json:"question_type"`
	Options      []QuestionOptionInput `json:"options"`
}

type QuestionOptionInput struct {
	Text string `json:"text"`
}
