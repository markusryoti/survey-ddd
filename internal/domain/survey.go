package domain

import (
	"errors"
	"fmt"
	"time"
)

type SurveyId string

type Survey struct {
	Id              SurveyId
	Title           string
	Description     *string
	MaxParticipants int
	AnswersReceived int
	EndTime         time.Time
	Questions       []Question
	SurveyStatus    SurveyStatus
	CreatedAt       time.Time
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
	return &Survey{
		Title:        title,
		Description:  description,
		SurveyStatus: Draft,
		CreatedAt:    time.Now(),
	}, nil
}

func (s *Survey) SetMaxParticipants(participants int) {
	s.MaxParticipants = participants
}

func (s *Survey) SetEndTime(endTime time.Time) error {
	if endTime.Before(time.Now()) {
		return errors.New("can't set end time that's in the past")
	}

	s.EndTime = endTime
	return nil
}

func (s *Survey) AddQuestion(question Question) {
	s.Questions = append(s.Questions, question)
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

	s.SurveyStatus = Released

	return nil
}

func (s *Survey) SubmissionReceived(receivedAt time.Time) error {
	if s.SurveyStatus == Draft {
		return errors.New("can't add a submission to a draft survey")
	}

	if s.SurveyStatus == Locked {
		return errors.New("can't add a submission to locked survey")
	}

	if s.AnswersReceived >= s.MaxParticipants {
		return fmt.Errorf("number of participants (%d) exceeded", s.MaxParticipants)
	}

	if s.EndTime.Before(receivedAt) {
		return errors.New("end time for survey has passed")
	}

	s.AnswersReceived++

	if s.AnswersReceived == s.MaxParticipants {
		s.SurveyStatus = Completed
	}

	return nil
}

func (s *Survey) Lock() {
	s.SurveyStatus = Locked
}

func (s Survey) Status() SurveyStatus {
	return s.SurveyStatus
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
