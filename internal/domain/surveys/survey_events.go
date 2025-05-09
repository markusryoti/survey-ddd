package surveys

import (
	"time"

	"github.com/markusryoti/survey-ddd/internal/core"
)

type SurveyCreated struct {
	Id           SurveyId
	Title        string
	Description  *string
	SurveyStatus SurveyStatus
	CreatedAt    time.Time
}

func (e SurveyCreated) AggregateId() core.AggregateId {
	return core.AggregateId(e.Id)
}

func (e SurveyCreated) Type() string {
	return "survey-created"
}

func (e SurveyCreated) Timestamp() time.Time {
	return e.CreatedAt
}

type QuestionAdded struct {
	Id        SurveyId
	Question  Question
	CreatedAt time.Time
}

func (e QuestionAdded) AggregateId() core.AggregateId {
	return core.AggregateId(e.Id)
}

func (e QuestionAdded) Type() string {
	return "question-added"
}

func (e QuestionAdded) Timestamp() time.Time {
	return e.CreatedAt
}

type MaxParticipantsChanged struct {
	Id              SurveyId
	MaxParticipants int
	CreatedAt       time.Time
}

func (e MaxParticipantsChanged) AggregateId() core.AggregateId {
	return core.AggregateId(e.Id)
}

func (e MaxParticipantsChanged) Type() string {
	return "max-participants-changed"
}

func (e MaxParticipantsChanged) Timestamp() time.Time {
	return e.CreatedAt
}

type SurveyEndTimeChanged struct {
	Id        SurveyId
	EndTime   time.Time
	CreatedAt time.Time
}

func (e SurveyEndTimeChanged) AggregateId() core.AggregateId {
	return core.AggregateId(e.Id)
}

func (e SurveyEndTimeChanged) Type() string {
	return "survey-endtime-changed"
}

func (e SurveyEndTimeChanged) Timestamp() time.Time {
	return e.CreatedAt
}

type SurveyReleased struct {
	Id        SurveyId
	CreatedAt time.Time
}

func (e SurveyReleased) AggregateId() core.AggregateId {
	return core.AggregateId(e.Id)
}

func (e SurveyReleased) Type() string {
	return "survey-released"
}

func (e SurveyReleased) Timestamp() time.Time {
	return e.CreatedAt
}

type SubmissionReceived struct {
	Id         SurveyId
	ReceivedAt time.Time
	CreatedAt  time.Time
}

func (e SubmissionReceived) AggregateId() core.AggregateId {
	return core.AggregateId(e.Id)
}

func (e SubmissionReceived) Type() string {
	return "submission-received"
}

func (e SubmissionReceived) Timestamp() time.Time {
	return e.CreatedAt
}

type SurveyCompleted struct {
	Id        SurveyId
	CreatedAt time.Time
}

func (e SurveyCompleted) AggregateId() core.AggregateId {
	return core.AggregateId(e.Id)
}

func (e SurveyCompleted) Type() string {
	return "survey-completed"
}

func (e SurveyCompleted) Timestamp() time.Time {
	return e.CreatedAt
}

type SurveyLocked struct {
	Id        SurveyId
	CreatedAt time.Time
}

func (e SurveyLocked) AggregateId() core.AggregateId {
	return core.AggregateId(e.Id)
}

func (e SurveyLocked) Type() string {
	return "survey-locked"
}

func (e SurveyLocked) Timestamp() time.Time {
	return e.CreatedAt
}
