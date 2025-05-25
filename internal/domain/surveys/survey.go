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

func (id *SurveyId) Scan(value any) error {
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
	SubmissionTimes []time.Time

	core.BaseAggregate
}

func NewSurvey(title string, description *string, tenantId string) (*Survey, error) {
	now := time.Now()

	if title == "" || tenantId == "" {
		return nil, errors.New("invalid survey")
	}

	survey := new(Survey)

	survey.addEvent(SurveyCreated{
		Id:           NewSurveyId(),
		Title:        title,
		Description:  description,
		TenantId:     tenantId,
		SurveyStatus: Draft,
		CreatedAt:    now,
	})

	return survey, nil
}

func (s Survey) ID() core.AggregateId {
	return core.AggregateId(s.Id)
}

func (s Survey) Name() string {
	return "survey"
}

func (s Survey) TableName() string {
	return "surveys"
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

func (s *Survey) SetMaxParticipants(participants int) error {
	if participants < 3 {
		return errors.New("min participants is three")
	}

	s.addEvent(MaxParticipantsChanged{
		Id:              s.Id,
		MaxParticipants: participants,
		CreatedAt:       time.Now(),
	})

	return nil
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
		s.TenantId = e.TenantId
		s.SurveyStatus = e.SurveyStatus
		s.SetCreatedAt(e.CreatedAt)
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
	s.AddDomainEvent(event)
	s.ApplyEvent(event)
}

func (s Survey) AnswersReceived() int {
	return len(s.SubmissionTimes)
}

func NewQuestion(title string, description string, options []string, allowMultiple bool) (Question, error) {
	if title == "" {
		return Question{}, errors.New("title cannot be empty")
	}

	if len(options) < 2 {
		return Question{}, errors.New("each options needs minimum of two options")
	}

	opts := make([]QuestionOption, 0)

	for _, opt := range options {
		o := NewQuestionOption(opt)
		opts = append(opts, *o)
	}

	var qt QuestionType

	if allowMultiple {
		qt = Multi
	} else {
		qt = Single
	}

	return Question{
		Id:              QuestionId(uuid.New().String()),
		Title:           title,
		Description:     &description,
		QuestionOptions: opts,
		QuestionType:    qt,
	}, nil
}

func NewQuestionOption(value string) *QuestionOption {
	return &QuestionOption{
		Id:    QuestionOptionId(uuid.New().String()),
		Value: value,
	}
}
