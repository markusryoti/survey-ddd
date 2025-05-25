package query

import (
	"context"

	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type QueryHandler struct {
	tx core.TransactionProvider
}

func NewQueryHandler(transactional core.TransactionProvider) *QueryHandler {
	return &QueryHandler{
		tx: transactional,
	}
}

func (q *QueryHandler) GetSurvey(ctx context.Context, id string) (surveys.Survey, error) {
	survey := new(surveys.Survey)

	surveyId, err := surveys.SurveyIdFromString(id)
	if err != nil {
		return *survey, err
	}

	err = q.tx.RunTransactional(ctx, func(repo core.Repository) error {
		return repo.Load(ctx, core.AggregateId(surveyId), survey)
	})

	return *survey, err
}
