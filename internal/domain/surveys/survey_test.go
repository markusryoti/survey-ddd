package surveys_test

import (
	"testing"
	"time"

	"github.com/markusryoti/survey-ddd/internal/domain/surveys"
	"github.com/stretchr/testify/assert"
)

func TestNewSurvey(t *testing.T) {
	t.Run("new survey is created and in draft state", func(t *testing.T) {
		title := "a title"
		description := "a description"
		survey, err := surveys.NewSurvey(title, &description, "tenant")

		assert.Nil(t, err)
		assert.Equal(t, title, survey.Title)
		assert.Equal(t, description, *survey.Description)
		assert.Equal(t, surveys.Draft, survey.SurveyStatus)
	})

	t.Run("can't add end time that is in the past", func(t *testing.T) {
		survey := newSurvey()
		err := survey.SetEndTime(now().Add(-1 * time.Minute))
		assert.NotNil(t, err)
	})
}

func TestReleaseAndLock(t *testing.T) {
	t.Run("can release survey", func(t *testing.T) {
		survey := newSurvey()
		survey.SetEndTime(now().Add(24 * time.Hour))
		survey.SetMaxParticipants(3)

		status := survey.Status()
		assert.Equal(t, surveys.Draft, status)

		err := survey.Release(now())
		assert.Nil(t, err)
		status = survey.Status()
		assert.Equal(t, surveys.Released, status)
	})

	t.Run("can't release if no max participants", func(t *testing.T) {
		survey := newSurvey()
		err := survey.Release(now())
		assert.NotNil(t, err)
	})

	t.Run("can lock survey", func(t *testing.T) {
		survey := newSurvey()
		survey.SetMaxParticipants(3)

		status := survey.Status()
		assert.Equal(t, surveys.Draft, status)

		err := survey.Release(now())
		assert.Nil(t, err)
		status = survey.Status()
		assert.Equal(t, surveys.Released, status)

		survey.Lock()
		status = survey.Status()
		assert.Equal(t, surveys.Locked, status)
	})
}

func TestSubmissions(t *testing.T) {
	t.Run("validate incorrect multioption answer to single question", func(t *testing.T) {
		survey := newSurvey()
		option1 := surveys.NewQuestionOption("option 1")
		option2 := surveys.NewQuestionOption("option 2")
		question, _ := surveys.NewQuestion("a guestion", "some stuff", []string{
			"option 1", "option 2",
		}, false)

		err := survey.ValidateResponse(question.Id, []surveys.QuestionOptionId{
			option1.Id, option2.Id,
		})
		assert.NotNil(t, err)
	})

	t.Run("can't create too many submissions", func(t *testing.T) {
		var err error

		survey := newSurvey()
		survey.SetMaxParticipants(3)
		survey.Release(now())

		err = survey.SubmissionReceived(now())
		assert.Nil(t, err)
		err = survey.SubmissionReceived(now())
		assert.Nil(t, err)
		err = survey.SubmissionReceived(now())
		assert.Nil(t, err)
		err = survey.SubmissionReceived(now())
		assert.NotNil(t, err)
	})

	t.Run("survey will be completed when max participants is achieved", func(t *testing.T) {
		survey := newSurvey()
		survey.SetMaxParticipants(3)
		survey.Release(now())

		_ = survey.SubmissionReceived(now())
		_ = survey.SubmissionReceived(now())
		_ = survey.SubmissionReceived(now())
		_ = survey.SubmissionReceived(now())

		status := survey.Status()
		assert.Equal(t, surveys.Completed, status)
	})

	t.Run("can't create a submission if survey is in draft state", func(t *testing.T) {
		survey := newSurvey()
		err := survey.SubmissionReceived(now())
		assert.NotNil(t, err)
	})

	t.Run("can't create a submission if survey is completed", func(t *testing.T) {
		survey := newSurvey()
		survey.SetMaxParticipants(3)
		survey.Release(now())

		_ = survey.SubmissionReceived(now())
		_ = survey.SubmissionReceived(now())
		_ = survey.SubmissionReceived(now())

		err := survey.SubmissionReceived(now())
		assert.NotNil(t, err)
	})

	t.Run("can't submit if survey is locked", func(t *testing.T) {
		survey := newSurvey()
		survey.Lock()
		err := survey.SubmissionReceived(now())
		assert.NotNil(t, err)
	})

	t.Run("can't submit if end time is in the past", func(t *testing.T) {
		survey := newSurvey()
		survey.SetMaxParticipants(3)

		survey.Release(now())

		survey.SetEndTime(now())
		err := survey.SubmissionReceived(now().Add(1 * time.Minute))
		assert.NotNil(t, err)
	})
}

func newSurvey() *surveys.Survey {
	description := "a description"
	survey, _ := surveys.NewSurvey("a title", &description, "tenant")
	survey.SetEndTime(now().Add(1 * time.Minute))
	return survey
}

func now() time.Time {
	return time.Now()
}
