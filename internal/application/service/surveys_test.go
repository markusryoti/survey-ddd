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

		repo := newMockRepo()
		transctional := newMockTransactionalProvider()

		description := "survey description"

		survey, err := surveys.NewSurvey("some title", &description, "tenant")
		assert.Nil(t, err)

		err = repo.Save(ctx, survey)
		assert.Nil(t, err)
		assert.NotEqual(t, "", survey.Id.String())

		srv := service.NewSurveyService(repo, transctional)

		err = srv.AddResponseToQuestion(ctx, service.ResponseToSurveyCmd{
			SurveyId: survey.Id.String(),
		})
		assert.Nil(t, err)

	})
}

type mockRepo struct {
}

func newMockRepo() *mockRepo {
	return &mockRepo{}
}

func (r *mockRepo) Load(ctx context.Context, id core.AggregateId, aggregate core.Aggregate) error {
	return nil
}

func (r *mockRepo) Save(ctx context.Context, aggregate core.Aggregate) error {

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
