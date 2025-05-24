package service

import (
	"context"
	"time"

	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type SurveyService struct {
	repo       core.Repository
	txProvider core.TransactionProvider
}

func NewSurveyService(
	repo core.Repository,
	txProvider core.TransactionProvider,
) *SurveyService {
	return &SurveyService{
		repo:       repo,
		txProvider: txProvider,
	}
}

type ResponseToSurveyCmd struct {
	SurveyId string
}

func (s *SurveyService) AddResponseToQuestion(ctx context.Context, cmd ResponseToSurveyCmd) error {
	err := s.txProvider.RunTransactional(ctx, func(repo core.Repository) error {
		surveyId, err := surveys.SurveyIdFromString(cmd.SurveyId)
		if err != nil {
			return err
		}

		response := surveys.NewSurveyResponse(surveyId)

		survey := new(surveys.Survey)

		err = s.repo.Load(ctx, core.AggregateId(surveyId), survey)
		if err != nil {
			return err
		}

		err = survey.SubmissionReceived(time.Now())
		if err != nil {
			return err
		}

		err = s.repo.Save(ctx, response)
		if err != nil {
			return err
		}

		err = s.repo.Save(ctx, survey)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
