package service

import (
	"context"
	"time"

	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type SurveyService struct {
	surveyRepo   core.Repository[*surveys.Survey]
	responseRepo core.Repository[*surveys.SurveyResponse]
	txProvider   core.TransactionProvider
}

func NewSurveyService(
	surveyRepo core.Repository[*surveys.Survey],
	responseRepo core.Repository[*surveys.SurveyResponse],
	txProvider core.TransactionProvider,
) *SurveyService {
	return &SurveyService{
		surveyRepo:   surveyRepo,
		responseRepo: responseRepo,
		txProvider:   txProvider,
	}
}

type ResponseToSurveyCmd struct {
	SurveyId string
}

func (s *SurveyService) AddResponseToQuestion(ctx context.Context, cmd ResponseToSurveyCmd) error {
	err := s.txProvider.RunTransactional(ctx, func(tx core.Transaction) error {
		surveyId, err := surveys.SurveyIdFromString(cmd.SurveyId)
		if err != nil {
			return err
		}

		response := surveys.NewSurveyResponse(surveyId)

		survey := new(surveys.Survey)

		err = s.surveyRepo.LoadWithTx(ctx, tx, core.AggregateId(surveyId), survey)
		if err != nil {
			return err
		}

		err = survey.SubmissionReceived(time.Now())
		if err != nil {
			return err
		}

		err = s.responseRepo.SaveWithTx(ctx, tx, response)
		if err != nil {
			return err
		}

		err = s.surveyRepo.SaveWithTx(ctx, tx, survey)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
