package service

import (
	"context"
	"database/sql"

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
}

func (s *SurveyService) AddResponseToQuestion(ctx context.Context, cmd ResponseToSurveyCmd) error {
	_, err := s.newTx(ctx)
	if err != nil {
		return err
	}

	return nil
}
