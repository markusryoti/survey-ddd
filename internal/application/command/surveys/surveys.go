package surveys

import (
	"context"
	"fmt"

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

func (h *SurveyCommandHandler) Handle(ctx context.Context, command core.Command) error {
	switch cmd := command.(type) {
	case surveys.CreateSurveyCommand:
		return h.handleCreateSurvey(ctx, cmd)
	case surveys.SetMaxParticipantsCommand:
		return h.handleSetMaxParticipants(ctx, cmd)
	default:
		return fmt.Errorf("unknown command type: %s", cmd.Type())
	}
}

func (h *SurveyCommandHandler) handleCreateSurvey(ctx context.Context, cmd surveys.CreateSurveyCommand) error {
	s, err := surveys.NewSurvey(cmd.Title, cmd.Description)
	if err != nil {
		return err
	}

	return h.repo.Save(ctx, s)
}

func (h *SurveyCommandHandler) handleSetMaxParticipants(ctx context.Context, cmd surveys.SetMaxParticipantsCommand) error {
	surveyId, err := surveys.SurveyIdFromString(cmd.SurveyId)

	survey, err := h.repo.Load(ctx, core.AggregateId(surveyId))
	if err != nil {
		return err
	}

	survey.SetMaxParticipants(cmd.MaxParticipants)

	return h.repo.Save(ctx, survey)
}
