package surveys_test

import (
	"testing"

	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
	"github.com/stretchr/testify/assert"
)

func TestNewSurveyResponse(t *testing.T) {
	t.Run("can't create multiple answers to single question", func(t *testing.T) {
		survey := newSurvey()
		question, _ := surveys.NewQuestion("a guestion", "some stuff", []string{
			"option 1", "option 2",
		}, false)

		response := surveys.NewSurveyResponse(survey.Id)
		err := response.AddResponseToQuestion(question.Id, []surveys.QuestionOptionId{
			question.QuestionOptions[0].Id})
		assert.Nil(t, err)

		err = response.AddResponseToQuestion(question.Id, []surveys.QuestionOptionId{
			question.QuestionOptions[0].Id,
		})
		assert.NotNil(t, err)
	})

	t.Run("can't create multiple answers to multi question", func(t *testing.T) {
		survey := newSurvey()
		question, _ := surveys.NewQuestion("a guestion", "some stuff", []string{
			"option1", "option2",
		}, true)

		response := surveys.NewSurveyResponse(survey.Id)
		err := response.AddResponseToQuestion(question.Id, []surveys.QuestionOptionId{
			question.QuestionOptions[0].Id,
		})
		assert.Nil(t, err)

		err = response.AddResponseToQuestion(question.Id, []surveys.QuestionOptionId{
			question.QuestionOptions[0].Id,
		})
		assert.NotNil(t, err)
	})
}
