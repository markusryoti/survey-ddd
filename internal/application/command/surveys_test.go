package command_test

import (
	"context"
	"testing"

	"github.com/markusryoti/survey-ddd/internal/application/command"
	"github.com/markusryoti/survey-ddd/internal/core"
	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
	"github.com/stretchr/testify/assert"
)

func TestCreateSurvey(t *testing.T) {
	t.Run("can create a survey", func(t *testing.T) {
		repo := newSurveyMockRepo[*surveys.Survey]()
		transctional := newMockTransactional()
		handler := command.NewSurveyCommandHandler[*surveys.Survey](repo, transctional)

		description := "survey description"

		survey, err := handler.HandleCreateSurvey(context.Background(), surveys.CreateSurveyCommand{
			Title:       "survey title",
			Description: &description,
		})

		assert.Nil(t, err)
		assert.Len(t, survey.GetUncommittedEvents(), 0)
	})
}

func TestSetMaxParticipants(t *testing.T) {
	t.Run("can create a survey", func(t *testing.T) {
		ctx := context.Background()
		repo := newSurveyMockRepo[*surveys.Survey]()
		transctional := newMockTransactional()
		handler := command.NewSurveyCommandHandler[*surveys.Survey](repo, transctional)

		description := "survey description"

		survey, err := handler.HandleCreateSurvey(ctx, surveys.CreateSurveyCommand{
			Title:       "survey title",
			Description: &description,
		})

		err = handler.HandleSetMaxParticipants(ctx, surveys.SetMaxParticipantsCommand{
			SurveyId:        survey.Id.String(),
			MaxParticipants: 3,
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

type mockTransactional struct {
}

func newMockTransactional() *mockTransactional {
	return &mockTransactional{}
}

func (t *mockTransactional) Begin(context.Context) error {
	return nil
}

func (t *mockTransactional) Commit() error {
	return nil
}

func (t *mockTransactional) Rollback() error {
	return nil
}

func (t *mockTransactional) RunTransactional(ctx context.Context, fn core.TransactionalSignature) error {
	return nil
}
