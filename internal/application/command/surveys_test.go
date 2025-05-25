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
		transctional := newMockTransactionalProvider()
		handler := command.NewCommandHandler(transctional)

		description := "survey description"

		survey, err := handler.CreateSurvey(context.Background(), surveys.CreateSurveyCommand{
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
		transctional := newMockTransactionalProvider()
		handler := command.NewCommandHandler(transctional)

		description := "survey description"

		survey, _ := handler.CreateSurvey(ctx, surveys.CreateSurveyCommand{
			Title:       "survey title",
			Description: &description,
		})

		err := handler.SetMaxParticipants(ctx, surveys.SetMaxParticipantsCommand{
			SurveyId:        survey.Id.String(),
			MaxParticipants: 3,
		})

		assert.Nil(t, err)
	})

	t.Run("can add a question", func(t *testing.T) {
		ctx := context.Background()
		transctional := newMockTransactionalProvider()
		handler := command.NewCommandHandler(transctional)

		title := "some title"
		description := "some description"

		survey, _ := handler.CreateSurvey(ctx, surveys.CreateSurveyCommand{
			Title:       title,
			Description: &description,
		})

		questionDescription := "question description"

		err := handler.AddQuestion(ctx, surveys.AddQuestionCommand{
			SurveyId: survey.Id.String(),
			Title:    "some question", Description: &questionDescription, AllowMultiple: false, QuestionOptions: []string{"a", "b"},
		})
		assert.Nil(t, err)
	})
}

type mockTransactionalProvider struct {
}

func newMockTransactionalProvider() *mockTransactionalProvider {
	return &mockTransactionalProvider{}
}

func (t *mockTransactionalProvider) RunTransactional(ctx context.Context, fn core.TransactionSignature) error {
	return nil
}
