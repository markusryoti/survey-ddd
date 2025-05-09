package surveys

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/markusryoti/survey-ddd/internal/core"
)

type SurveyResponseId uuid.UUID

type SurveyResponse struct {
	Id                SurveyResponseId
	SurveyId          SurveyId
	NumberOfQuestions int
	Responses         []QuestionResponse
	CreatedAt         time.Time
	Status            ResponseStatus

	version           int
	uncommittedEvents []core.Event
}

type ResponseStatus string

const (
	ResponseStatusDraft     ResponseStatus = "draft"
	ResponseStatusSubmitted ResponseStatus = "submitted"
)

type QuestionResponse struct {
	QuestionId QuestionId
	Choices    []QuestionOptionId
}

func NewSurveyResponse(id SurveyId, numQuestions int) *SurveyResponse {
	now := time.Now()

	response := &SurveyResponse{}

	response.addEvent(SurveyResponseCreated{
		SurveyId:          id,
		NumberOfQuestions: numQuestions,
		CreatedAt:         now,
	})

	return response
}

func (s *SurveyResponse) AddResponseToQuestion(question QuestionId, options []QuestionOptionId) error {
	for _, q := range s.Responses {
		if q.QuestionId == question {
			return errors.New("not allowed to answer multiple times")
		}
	}

	s.addEvent(QuestionAnswered{
		QuestionId: question,
		Choices:    options,
	})

	return nil
}

func (s *SurveyResponse) Submit() {
	s.addEvent(ResponseSubmitted{
		SurveyId:  s.SurveyId,
		CreatedAt: time.Now(),
	})
}

func (s *SurveyResponse) ApplyEvent(event core.Event) {
	switch e := event.(type) {
	case SurveyResponseCreated:
		s.SurveyId = e.SurveyId
		s.NumberOfQuestions = e.NumberOfQuestions
		s.CreatedAt = e.CreatedAt
	case QuestionAnswered:
		s.Responses = append(s.Responses, QuestionResponse{
			QuestionId: e.QuestionId,
			Choices:    e.Choices,
		})
	case ResponseSubmitted:
		s.Status = ResponseStatusSubmitted
	default:
		panic(fmt.Sprintf("unknown event: %+v", e))
	}
}

func (s *SurveyResponse) addEvent(event core.Event) {
	s.uncommittedEvents = append(s.uncommittedEvents, event)
	s.ApplyEvent(event)
}
