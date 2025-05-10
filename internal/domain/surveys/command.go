package surveys

type CreateSurveyCommand struct {
	Title       string
	Description *string
}

func (c CreateSurveyCommand) Type() string {
	return "CreateSurveyCommand"
}

type SetMaxParticipantsCommand struct {
	SurveyId        string
	MaxParticipants int
}

func (c SetMaxParticipantsCommand) Type() string {
	return "SetMaxParticipantsCommand"
}
