package surveys

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/markusryoti/survey-ddd/internal/core"
)

type SurveyId core.AggregateId

func NewSurveyId() SurveyId {
	return SurveyId(core.NewAggregateId())
}

func (s SurveyId) String() string {
	return core.AggregateId(s).String()
}

func (id SurveyId) MarshalJSON() ([]byte, error) {
	return core.AggregateId(id).MarshalJSON()
}

func (id *SurveyId) UnmarshalJSON(data []byte) error {
	return (*core.AggregateId)(id).UnmarshalJSON(data)
}

func (id SurveyId) Value() (driver.Value, error) {
	return core.AggregateId(id).Value()
}

func (id *SurveyId) Scan(value interface{}) error {
	return (*core.AggregateId)(id).Scan(value)
}

func SurveyIdFromString(s string) (SurveyId, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return SurveyId{}, err
	}

	return SurveyId(core.AggregateId(id)), nil
}

type Survey struct {
	Id              SurveyId
	Title           string
	Description     *string
	MaxParticipants int
	EndTime         time.Time
	Questions       []Question
	SurveyStatus    SurveyStatus
	TenantId        string
	TimeCreated     time.Time
	SubmissionTimes []time.Time

	version           int
	uncommittedEvents []core.DomainEvent
}

func (s Survey) ID() core.AggregateId {
	return core.AggregateId(s.Id)
}

func (s Survey) GetUncommittedEvents() []core.DomainEvent {
	return s.uncommittedEvents
}

func (s *Survey) ClearUncommittedEvents() {
	s.uncommittedEvents = make([]core.DomainEvent, 0)
}

func (s *Survey) SetVersion(version int) {
	s.version = version
}

func (s *Survey) SetCreatedAt(t time.Time) {
	s.TimeCreated = t
}

func (s Survey) Version() int {
	return s.version
}

func (s Survey) CreatedAt() time.Time {
	return s.TimeCreated
}

type SurveyStatus string

const (
	Draft     SurveyStatus = "draft"
	Released  SurveyStatus = "released"
	Locked    SurveyStatus = "locked"
	Completed SurveyStatus = "completed"
)

type QuestionId string

type Question struct {
	Id              QuestionId
	Title           string
	Description     *string
	QuestionType    QuestionType
	QuestionOptions []QuestionOption
}

type QuestionType string

const (
	Single QuestionType = "single"
	Multi  QuestionType = "multi"
)

func NewQuestionType(questionType string) (QuestionType, error) {
	var qt QuestionType
	switch questionType {
	case "single":
		qt = QuestionType("single")
	case "multi":
		qt = QuestionType("multi")
	default:
		return "", fmt.Errorf("invalid question type: %s", questionType)
	}

	return qt, nil
}

type QuestionOptionId string

type QuestionOption struct {
	Id    QuestionOptionId
	Value string
}

func NewSurvey(title string, description *string) (*Survey, error) {
	now := time.Now()

	survey := &Survey{}

	survey.addEvent(SurveyCreated{
		Id:           NewSurveyId(),
		Title:        title,
		Description:  description,
		SurveyStatus: Draft,
		CreatedAt:    now,
	})

	return survey, nil
}

func (s *Survey) SetMaxParticipants(participants int) {
	s.addEvent(MaxParticipantsChanged{
		Id:              s.Id,
		MaxParticipants: participants,
		CreatedAt:       time.Now(),
	})
}

func (s *Survey) SetEndTime(endTime time.Time) error {
	if endTime.Before(time.Now()) {
		return errors.New("can't set end time that's in the past")
	}

	s.addEvent(SurveyEndTimeChanged{
		Id:        s.Id,
		EndTime:   endTime,
		CreatedAt: time.Now(),
	})

	return nil
}

func (s *Survey) AddQuestion(question Question) {
	s.addEvent(QuestionAdded{
		Id:        s.Id,
		Question:  question,
		CreatedAt: time.Now(),
	})
}

func (s *Survey) Release(now time.Time) error {
	if s.MaxParticipants == 0 {
		return errors.New("can't release without number of participants")
	}

	if s.EndTime.IsZero() {
		return errors.New("can't release without end time")
	}

	if s.EndTime.Before(now) {
		return errors.New("can't release if end time is in the past")
	}

	s.addEvent(SurveyReleased{
		Id:        s.Id,
		CreatedAt: now,
	})

	return nil
}

func (s Survey) ValidateResponse(question QuestionId, options []QuestionOptionId) error {
	q, err := s.getQuestion(question)
	if err != nil {
		return err
	}

	if q.QuestionType == Single && len(options) > 1 {
		return errors.New("not allowed to answer with multiple options")
	}

	return nil
}

func (s *Survey) getQuestion(id QuestionId) (Question, error) {
	for _, q := range s.Questions {
		if q.Id == id {
			return q, nil
		}
	}

	return Question{}, errors.New("question not found")
}

func (s *Survey) SubmissionReceived(receivedAt time.Time) error {
	if s.SurveyStatus == Draft {
		return errors.New("can't add a submission to a draft survey")
	}

	if s.SurveyStatus == Locked {
		return errors.New("can't add a submission to locked survey")
	}

	if s.AnswersReceived() >= s.MaxParticipants {
		return fmt.Errorf("number of participants (%d) exceeded", s.MaxParticipants)
	}

	if s.EndTime.Before(receivedAt) {
		return errors.New("end time for survey has passed")
	}

	s.addEvent(SubmissionReceived{
		Id:         s.Id,
		ReceivedAt: receivedAt,
		CreatedAt:  time.Now(),
	})

	if s.AnswersReceived() == s.MaxParticipants {
		s.addEvent(SurveyCompleted{
			Id:        s.Id,
			CreatedAt: time.Now(),
		})
	}

	return nil
}

func (s *Survey) Lock() {
	s.addEvent(SurveyLocked{
		Id:        s.Id,
		CreatedAt: time.Now(),
	})
}

func (s Survey) Status() SurveyStatus {
	return s.SurveyStatus
}

func (s *Survey) ApplyEvent(event core.DomainEvent) {
	switch e := event.(type) {
	case SurveyCreated:
		s.Id = e.Id
		s.Title = e.Title
		s.Description = e.Description
		s.SurveyStatus = e.SurveyStatus
		s.TimeCreated = e.CreatedAt
	case QuestionAdded:
		s.Questions = append(s.Questions, e.Question)
	case MaxParticipantsChanged:
		s.MaxParticipants = e.MaxParticipants
	case SurveyEndTimeChanged:
		s.EndTime = e.EndTime
	case SurveyReleased:
		s.SurveyStatus = Released
	case SubmissionReceived:
		s.SubmissionTimes = append(s.SubmissionTimes, e.ReceivedAt)
	case SurveyCompleted:
		s.SurveyStatus = Completed
	case SurveyLocked:
		s.SurveyStatus = Locked
	default:
		panic(fmt.Sprintf("unknown event: %+v", e))
	}
}

func (s *Survey) addEvent(event core.DomainEvent) {
	s.uncommittedEvents = append(s.uncommittedEvents, event)
	s.ApplyEvent(event)
}

func (s Survey) AnswersReceived() int {
	return len(s.SubmissionTimes)
}

func NewQuestion(title string, description string, options []QuestionOption, allowMultiple bool) *Question {
	var qt QuestionType

	if allowMultiple {
		qt = Multi
	} else {
		qt = Single
	}

	return &Question{
		Title:           title,
		Description:     &description,
		QuestionOptions: options,
		QuestionType:    qt,
	}
}

func NewQuestionOption(value string) *QuestionOption {
	return &QuestionOption{
		Value: value,
	}
}
