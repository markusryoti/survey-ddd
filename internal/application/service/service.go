package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/markusryoti/survey-ddd/internal/adapters/postgres"
	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
)

type SurveyService struct {
	surveyRepo           *postgres.PostgresRepository[*surveys.Survey]
	responseRepo         *postgres.PostgresRepository[*surveys.SurveyResponse]
	domainEventDipatcher *core.EventDispatcher
	db                   *sql.DB
}

func NewSurveyService(
	surveyRepo *postgres.PostgresRepository[*surveys.Survey],
	responseRepo *postgres.PostgresRepository[*surveys.SurveyResponse],
	domainEventDipatcher *core.EventDispatcher,
	db *sql.DB,
) *SurveyService {
	return &SurveyService{
		surveyRepo:           surveyRepo,
		responseRepo:         responseRepo,
		domainEventDipatcher: domainEventDipatcher,
		db:                   db,
	}
}

func (s *SurveyService) newTx(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, &sql.TxOptions{})
}

type ResponseToSurveyCmd struct {
	SurveyId string
}

func (s *SurveyService) AddResponseToQuestion(ctx context.Context, cmd ResponseToSurveyCmd) error {
	var err error

	tx, err := s.newTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				panic("rollback failed")
			}
		}
	}()

	var survey surveys.Survey
	surveyId, err := surveys.SurveyIdFromString(cmd.SurveyId)
	if err != nil {
		return err
	}

	response := surveys.NewSurveyResponse(surveyId)

	err = s.surveyRepo.Load(ctx, core.AggregateId(surveyId), &survey)
	if err != nil {
		return err
	}

	err = survey.SubmissionReceived(time.Now())
	if err != nil {
		return err
	}

	err = s.responseRepo.Save(ctx, response)
	if err != nil {
		return err
	}

	err = s.surveyRepo.Save(ctx, &survey)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
