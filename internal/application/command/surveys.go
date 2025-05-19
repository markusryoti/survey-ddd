package command

import (
	"context"
	"time"

	"github.com/markusryoti/survey-ddd/internal/adapters/postgres"
	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type SurveyCommandHandler struct {
	repo                 core.Repository[*surveys.Survey]
	domainEventDipatcher *core.EventDispatcher
}

func NewSurveyCommandHandler[T core.Aggregate](
	repo *postgres.PostgresRepository[*surveys.Survey],
) *SurveyCommandHandler {
	return &SurveyCommandHandler{
		repo:                 repo,
		domainEventDipatcher: core.NewEventDispatcher(),
	}
}

func (h *SurveyCommandHandler) HandleCreateSurvey(ctx context.Context, cmd surveys.CreateSurveyCommand) (*surveys.Survey, error) {
	survey, err := surveys.NewSurvey(cmd.Title, cmd.Description)
	if err != nil {
		return survey, err
	}

	err = h.repo.Save(ctx, survey)
	if err != nil {
		return survey, err
	}

	h.domainEventDipatcher.Dispatch(ctx, survey.GetUncommittedEvents()...)

	survey.ClearUncommittedEvents()

	return survey, err
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
