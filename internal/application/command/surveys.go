package command

import (
	"context"

	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type CommandHandler struct {
	tx core.TransactionProvider
}

func NewCommandHandler(
	txProvider core.TransactionProvider,
) *CommandHandler {
	return &CommandHandler{
		tx: txProvider,
	}
}

func (h *CommandHandler) CreateSurvey(ctx context.Context, cmd surveys.CreateSurveyCommand) (*surveys.Survey, error) {
	var err error

	survey := new(surveys.Survey)

	err = h.tx.RunTransactional(ctx, func(repo core.Repository) error {
		defer survey.ClearUncommittedEvents()

		survey, err = surveys.NewSurvey(cmd.Title, cmd.Description, cmd.TenantId)
		if err != nil {
			return err
		}

		err = repo.Save(ctx, survey)
		if err != nil {
			return err
		}

		return err
	})

	return survey, err
}

func (h *CommandHandler) SetMaxParticipants(ctx context.Context, cmd surveys.SetMaxParticipantsCommand) error {
	surveyId, err := surveys.SurveyIdFromString(cmd.SurveyId)
	if err != nil {
		return err
	}

	return h.tx.RunTransactional(ctx, func(repo core.Repository) error {
		survey := new(surveys.Survey)

		defer survey.ClearUncommittedEvents()

		err = repo.Load(ctx, core.AggregateId(surveyId), survey)
		if err != nil {
			return err
		}

		err = survey.SetMaxParticipants(cmd.MaxParticipants)
		if err != nil {
			return err
		}

		err = repo.Save(ctx, survey)
		if err != nil {
			return err
		}

		return nil
	})
}

func (h *CommandHandler) AddQuestion(ctx context.Context, cmd surveys.AddQuestionCommand) error {
	surveyId, err := surveys.SurveyIdFromString(cmd.SurveyId)
	if err != nil {
		return err
	}

	q, err := surveys.NewQuestion(cmd.Title, *cmd.Description, cmd.QuestionOptions, cmd.AllowMultiple)
	if err != nil {
		return err
	}

	return h.tx.RunTransactional(ctx, func(repo core.Repository) error {
		survey := new(surveys.Survey)

		defer survey.ClearUncommittedEvents()

		err := repo.Load(ctx, core.AggregateId(surveyId), survey)
		if err != nil {
			return err
		}

		survey.AddQuestion(q)

		err = repo.Save(ctx, survey)
		if err != nil {
			return err
		}

		return nil
	})
}
