package domain_test

import (
	"testing"

	"github.com/markusryoti/survey-ddd/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewSurveyResponse(t *testing.T) {
	t.Run("can't create multioption answer to single question", func(t *testing.T) {
		survey := newSurvey()
		option1 := domain.NewQuestionOption("option 1")
		option2 := domain.NewQuestionOption("option 2")
		question := domain.NewQuestion("a guestion", "some stuff", []domain.QuestionOption{
			*option1, *option2,
		}, false)

		response := domain.NewSurveyResponse(survey.Id)
		err := response.AddResponseToQuestion(*question, []domain.QuestionOption{
			*option1, *option2,
		})
		assert.NotNil(t, err)
	})

	t.Run("can create multioption answer to multioption question", func(t *testing.T) {
		survey := newSurvey()
		option1 := domain.NewQuestionOption("option 1")
		option2 := domain.NewQuestionOption("option 2")
		question := domain.NewQuestion("a guestion", "some stuff", []domain.QuestionOption{
			*option1, *option2,
		}, true)

		response := domain.NewSurveyResponse(survey.Id)
		err := response.AddResponseToQuestion(*question, []domain.QuestionOption{
			*option1, *option2,
		})
		assert.Nil(t, err)
	})

	t.Run("can't create multiple answers to single question", func(t *testing.T) {
		survey := newSurvey()
		option1 := domain.NewQuestionOption("option 1")
		option2 := domain.NewQuestionOption("option 2")
		question := domain.NewQuestion("a guestion", "some stuff", []domain.QuestionOption{
			*option1, *option2,
		}, false)

		response := domain.NewSurveyResponse(survey.Id)
		err := response.AddResponseToQuestion(*question, []domain.QuestionOption{
			*option1,
		})
		assert.Nil(t, err)

		err = response.AddResponseToQuestion(*question, []domain.QuestionOption{
			*option1,
		})
		assert.NotNil(t, err)
	})

	t.Run("can't create multiple answers to multi question", func(t *testing.T) {
		survey := newSurvey()
		option1 := domain.NewQuestionOption("option 1")
		option2 := domain.NewQuestionOption("option 2")
		question := domain.NewQuestion("a guestion", "some stuff", []domain.QuestionOption{
			*option1, *option2,
		}, true)

		response := domain.NewSurveyResponse(survey.Id)
		err := response.AddResponseToQuestion(*question, []domain.QuestionOption{
			*option1,
		})
		assert.Nil(t, err)

		err = response.AddResponseToQuestion(*question, []domain.QuestionOption{
			*option1,
		})
		assert.NotNil(t, err)
	})
}
