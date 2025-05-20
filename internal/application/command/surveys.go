package command

import (
	"context"
	"time"

	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type SurveyCommandHandler interface {
	HandleCreateSurvey(ctx context.Context, cmd surveys.CreateSurveyCommand) (*surveys.Survey, error)
	HandleSetMaxParticipants(ctx context.Context, cmd surveys.SetMaxParticipantsCommand) error
}

type SurveyPostgresCommandHandler struct {
	repo                 core.Repository[*surveys.Survey]
	domainEventDipatcher *core.EventDispatcher
	txProvider           core.TransactionProvider
}

func NewSurveyCommandHandler[T core.Aggregate](
	repo core.Repository[*surveys.Survey],
	txProvider core.TransactionProvider,
) *SurveyPostgresCommandHandler {
	return &SurveyPostgresCommandHandler{
		repo:                 repo,
		domainEventDipatcher: core.NewEventDispatcher(),
		txProvider:           txProvider,
	}
}

func (h *SurveyPostgresCommandHandler) HandleCreateSurvey(ctx context.Context, cmd surveys.CreateSurveyCommand) (*surveys.Survey, error) {
	var err error

	survey := new(surveys.Survey)

	err = h.txProvider.RunTransactional(ctx, func(tx core.Transaction) error {
		survey, err = surveys.NewSurvey(cmd.Title, cmd.Description)
		if err != nil {
			return err
		}

		err = h.repo.SaveWithTx(ctx, tx, survey)
		if err != nil {
			return err
		}

		h.domainEventDipatcher.Dispatch(ctx, survey.GetUncommittedEvents()...)

		survey.ClearUncommittedEvents()

		return err
	})

	return survey, err
}

func (h *SurveyPostgresCommandHandler) HandleSetMaxParticipants(ctx context.Context, cmd surveys.SetMaxParticipantsCommand) error {
	surveyId, err := surveys.SurveyIdFromString(cmd.SurveyId)

	survey := new(surveys.Survey)

	err = h.repo.Load(ctx, core.AggregateId(surveyId), survey)
	if err != nil {
		return err
	}

	err = survey.SetMaxParticipants(cmd.MaxParticipants)
	if err != nil {
		return err
	}

	err = h.repo.Save(ctx, survey)
	if err != nil {
		return err
	}

	h.domainEventDipatcher.Dispatch(ctx, survey.GetUncommittedEvents()...)

	survey.ClearUncommittedEvents()

	return nil
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
