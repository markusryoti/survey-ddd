package command

import (
	"context"

	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type SurveyCommandHandler interface {
	HandleCreateSurvey(ctx context.Context, cmd surveys.CreateSurveyCommand) (*surveys.Survey, error)
	HandleSetMaxParticipants(ctx context.Context, cmd surveys.SetMaxParticipantsCommand) error
}

type SurveyCmdHandler struct {
	repo       core.Repository[*surveys.Survey]
	txProvider core.TransactionProvider
}

func NewSurveyCommandHandler[T core.Aggregate](
	repo core.Repository[*surveys.Survey],
	txProvider core.TransactionProvider,
) *SurveyCmdHandler {
	return &SurveyCmdHandler{
		repo:       repo,
		txProvider: txProvider,
	}
}

func (h *SurveyCmdHandler) HandleCreateSurvey(ctx context.Context, cmd surveys.CreateSurveyCommand) (*surveys.Survey, error) {
	var err error

	survey := new(surveys.Survey)

	err = h.txProvider.RunTransactional(ctx, func(tx core.Transaction) error {
		survey, err = surveys.NewSurvey(cmd.Title, cmd.Description, cmd.TenantId)
		if err != nil {
			return err
		}

		err = h.repo.SaveWithTx(ctx, tx, survey)
		if err != nil {
			return err
		}

		survey.ClearUncommittedEvents()

		return err
	})

	return survey, err
}

func (h *SurveyCmdHandler) HandleSetMaxParticipants(ctx context.Context, cmd surveys.SetMaxParticipantsCommand) error {
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

	survey.ClearUncommittedEvents()

	return nil
}

func (h *SurveyCmdHandler) AddQuestion(ctx context.Context, cmd surveys.AddQuestionCommand) error {
	surveyId, err := surveys.SurveyIdFromString(cmd.SurveyId)
	if err != nil {
		return err
	}

	q, err := surveys.NewQuestion(cmd.Title, *cmd.Description, cmd.QuestionOptions, cmd.AllowMultiple)
	if err != nil {
		return err
	}

	err = h.txProvider.RunTransactional(ctx, func(tx core.Transaction) error {
		survey := new(surveys.Survey)

		err := h.repo.LoadWithTx(ctx, tx, core.AggregateId(surveyId), survey)
		if err != nil {
			return err
		}

		survey.AddQuestion(*q)

		err = h.repo.SaveWithTx(ctx, tx, survey)

		survey.ClearUncommittedEvents()

		return err
	})

	return err
}
