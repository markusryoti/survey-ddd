package surveys_test

import (
	"testing"

	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
	"github.com/stretchr/testify/assert"
)

func TestNewSurveyResponse(t *testing.T) {
	t.Run("can't create multiple answers to single question", func(t *testing.T) {
		survey := newSurvey()
		option1 := surveys.NewQuestionOption("option 1")
		option2 := surveys.NewQuestionOption("option 2")
		question := surveys.NewQuestion("a guestion", "some stuff", []surveys.QuestionOption{
			*option1, *option2,
		}, false)

		response := surveys.NewSurveyResponse(survey.Id, 1)
		err := response.AddResponseToQuestion(question.Id,
			[]surveys.QuestionOptionId{option1.Id},
		)
		assert.Nil(t, err)

		err = response.AddResponseToQuestion(question.Id, []surveys.QuestionOptionId{
			option1.Id,
		})
		assert.NotNil(t, err)
	})

	t.Run("can't create multiple answers to multi question", func(t *testing.T) {
		survey := newSurvey()
		option1 := surveys.NewQuestionOption("option 1")
		option2 := surveys.NewQuestionOption("option 2")
		question := surveys.NewQuestion("a guestion", "some stuff", []surveys.QuestionOption{
			*option1, *option2,
		}, true)

		response := surveys.NewSurveyResponse(survey.Id, 1)
		err := response.AddResponseToQuestion(question.Id, []surveys.QuestionOptionId{
			option1.Id,
		})
		assert.Nil(t, err)

		err = response.AddResponseToQuestion(question.Id, []surveys.QuestionOptionId{
			option1.Id,
		})
		assert.NotNil(t, err)
	})
}
