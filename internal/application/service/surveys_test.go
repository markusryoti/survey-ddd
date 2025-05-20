package service_test

import (
	"context"
	"testing"

	"github.com/markusryoti/survey-ddd/internal/application/service"
	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
	"github.com/stretchr/testify/assert"
)

func TestSubmitResponse(t *testing.T) {
	t.Run("can submit a response to question", func(t *testing.T) {
		ctx := context.Background()

		surveyRepo := newSurveyMockRepo[*surveys.Survey]()
		responseRepo := newSurveyMockRepo[*surveys.SurveyResponse]()
		transctional := newMockTransactionalProvider()

		description := "survey description"

		survey, err := surveys.NewSurvey("some title", &description)
		assert.Nil(t, err)

		err = surveyRepo.Save(ctx, survey)
		assert.Nil(t, err)
		assert.NotEqual(t, "", survey.Id.String())

		srv := service.NewSurveyService(surveyRepo, responseRepo, nil, transctional)

		err = srv.AddResponseToQuestion(ctx, service.ResponseToSurveyCmd{
			SurveyId: survey.Id.String(),
		})
		assert.Nil(t, err)

	})
}

type mockRepo[T core.Aggregate] struct {
}

func newSurveyMockRepo[T core.Aggregate]() *mockRepo[T] {
	return &mockRepo[T]{}
}

func (r *mockRepo[T]) Load(ctx context.Context, id core.AggregateId, aggregate T) error {
	return nil
}

func (r *mockRepo[T]) Save(ctx context.Context, aggregate T) error {

	return nil
}

func (r *mockRepo[T]) SaveWithTx(ctx context.Context, tx core.Transaction, aggregate T) error {
	return nil
}

func (r *mockRepo[T]) LoadWithTx(ctx context.Context, tx core.Transaction, id core.AggregateId, aggregate T) error {
	return nil
}

type mockTransactionalProvider struct {
}

func newMockTransactionalProvider() *mockTransactionalProvider {
	return &mockTransactionalProvider{}
}

func (t *mockTransactionalProvider) RunTransactional(ctx context.Context, fn core.TransactionalSignature) error {
	return nil
}
