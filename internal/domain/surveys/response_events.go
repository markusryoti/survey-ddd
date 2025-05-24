package surveys

import (
	"time"

	"github.com/markusryoti/survey-ddd/internal/core"
)

type SurveyResponseCreated struct {
	Id                SurveyResponseId
	SurveyId          SurveyId
	NumberOfQuestions int
	CreatedAt         time.Time
}

func (e SurveyResponseCreated) AggregateId() core.AggregateId {
	return core.AggregateId(e.Id)
}

func (e SurveyResponseCreated) Type() string {
	return "survey-response-created"
}

func (e SurveyResponseCreated) OccurredAt() time.Time {
	return e.CreatedAt
}

type QuestionAnswered struct {
	Id         SurveyResponseId
	QuestionId QuestionId
	Choices    []QuestionOptionId
	CreatedAt  time.Time
}

func (e QuestionAnswered) AggregateId() core.AggregateId {
	return core.AggregateId(e.Id)
}

func (e QuestionAnswered) Type() string {
	return "question-answered"
}

func (e QuestionAnswered) OccurredAt() time.Time {
	return e.CreatedAt
}

type ResponseSubmitted struct {
	Id        SurveyResponseId
	SurveyId  SurveyId
	CreatedAt time.Time
}

func (e ResponseSubmitted) AggregateId() core.AggregateId {
	return core.AggregateId(e.Id)
}

func (e ResponseSubmitted) Type() string {
	return "response-submitted"
}

func (e ResponseSubmitted) OccurredAt() time.Time {
	return e.CreatedAt
}
