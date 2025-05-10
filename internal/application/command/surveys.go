package command

import (
	"context"
	"time"

	"github.com/markusryoti/survey-ddd/internal/adapters/postgres"
	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type SurveyCommandHandler struct {
	repo *postgres.PostgresRepository[*surveys.Survey]
}

func NewSurveyCommandHandler(
	repo *postgres.PostgresRepository[*surveys.Survey],
) *SurveyCommandHandler {
	return &SurveyCommandHandler{
		repo: repo,
	}
}

func (h *SurveyCommandHandler) HandleCreateSurvey(ctx context.Context, cmd surveys.CreateSurveyCommand) (*surveys.Survey, error) {
	s, err := surveys.NewSurvey(cmd.Title, cmd.Description)
	if err != nil {
		return s, err
	}

	err = h.repo.Save(ctx, s)
	return s, err
}

func (h *SurveyCommandHandler) HandleSetMaxParticipants(ctx context.Context, cmd surveys.SetMaxParticipantsCommand) error {
	surveyId, err := surveys.SurveyIdFromString(cmd.SurveyId)

	survey := &surveys.Survey{}

	err = h.repo.Load(ctx, core.AggregateId(surveyId), survey)
	if err != nil {
		return err
	}

	survey.SetMaxParticipants(cmd.MaxParticipants)

	return h.repo.Save(ctx, survey)
}

type CreateSurveyCommand struct {
	Title           string    `json:"title"`
	Description     *string   `json:"description"`
	MaxParticipants int       `json:"max_participants"`
	EndTime         time.Time `json:"end_time"`
	TenantID        string    `json:"tenant_id"`
}

func (c CreateSurveyCommand) Type() string {
	return "CreateSurveyCommand"
}
