package surveys

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/markusryoti/survey-ddd/internal/core"
)

type SurveyResponseId core.AggregateId

func (s SurveyResponseId) String() string {
	return core.AggregateId(s).String()
}

func (id SurveyResponseId) MarshalJSON() ([]byte, error) {
	return core.AggregateId(id).MarshalJSON()
}

func (id *SurveyResponseId) UnmarshalJSON(data []byte) error {
	return (*core.AggregateId)(id).UnmarshalJSON(data)
}

func (id SurveyResponseId) Value() (driver.Value, error) {
	return core.AggregateId(id).Value()
}

func (id *SurveyResponseId) Scan(value interface{}) error {
	return (*core.AggregateId)(id).Scan(value)
}

func SurveyResponseIdFromString(s string) (SurveyId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return SurveyId{}, err
	}

	return SurveyId(core.AggregateId(id)), nil
}

type SurveyResponse struct {
	Id                SurveyResponseId
	SurveyId          SurveyId
	NumberOfQuestions int
	Responses         []QuestionResponse
	TimeCreated       time.Time
	Status            ResponseStatus

	version           int
	uncommittedEvents []core.DomainEvent
}

func (s SurveyResponse) ID() core.AggregateId {
	return core.AggregateId(s.Id)
}

func (s SurveyResponse) GetUncommittedEvents() []core.DomainEvent {
	return s.uncommittedEvents
}

func (s *SurveyResponse) ClearUncommittedEvents() {
	s.uncommittedEvents = make([]core.DomainEvent, 0)
}

func (s *SurveyResponse) SetVersion(version int) {
	s.version = version
}

func (s *SurveyResponse) SetCreatedAt(t time.Time) {
	s.TimeCreated = t
}

func (s SurveyResponse) Version() int {
	return s.version
}

func (s SurveyResponse) CreatedAt() time.Time {
	return s.TimeCreated
}

func (s SurveyResponse) Name() string {
	return "survey-response"
}

func (s SurveyResponse) TableName() string {
	return "survey_responses"
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

func NewSurveyResponse(id SurveyId) *SurveyResponse {
	now := time.Now()

	response := &SurveyResponse{}

	response.addEvent(SurveyResponseCreated{
		SurveyId:  id,
		CreatedAt: now,
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

func (s *SurveyResponse) ApplyEvent(event core.DomainEvent) {
	switch e := event.(type) {
	case SurveyResponseCreated:
		s.SurveyId = e.SurveyId
		s.NumberOfQuestions = e.NumberOfQuestions
		s.TimeCreated = e.CreatedAt
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

func (s *SurveyResponse) addEvent(event core.DomainEvent) {
	s.uncommittedEvents = append(s.uncommittedEvents, event)
	s.ApplyEvent(event)
}
