package domain

import (
	"errors"
	"time"
)

type SurveyResponseId string

type SurveyResponse struct {
	Id        SurveyResponseId
	SurveyId  SurveyId
	Responses []QuestionResponse
	CreatedAt time.Time
}

type QuestionResponse struct {
	QuestionId QuestionId
	Choices    []QuestionOptionId
}

func NewSurveyResponse(id SurveyId) *SurveyResponse {
	return &SurveyResponse{
		SurveyId:  id,
		Responses: make([]QuestionResponse, 0),
		CreatedAt: time.Now(),
	}
}

func (s *SurveyResponse) AddResponseToQuestion(question Question, options []QuestionOption) error {
	if question.QuestionType == Single && len(options) > 1 {
		return errors.New("not allowed to answer with multiple options")
	}

	for _, q := range s.Responses {
		if q.QuestionId == question.Id {
			return errors.New("not allowed to answer multiple times")
		}
	}

	choices := make([]QuestionOptionId, 0)
	for _, c := range options {
		choices = append(choices, c.Id)
	}

	s.Responses = append(s.Responses, QuestionResponse{
		QuestionId: question.Id,
		Choices:    choices,
	})

	return nil
}
