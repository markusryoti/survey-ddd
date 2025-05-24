package surveys

type CreateSurveyCommand struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	TenantId    string  `json:"tenantId"`
}

type SetMaxParticipantsCommand struct {
	SurveyId        string `json:"surveyId"`
	MaxParticipants int    `json:"maxParticipants"`
}

type AddQuestionCommand struct {
	SurveyId        string   `json:"surveyId"`
	Title           string   `json:"title"`
	Description     *string  `json:"description"`
	AllowMultiple   bool     `json:"allowMultiple"`
	QuestionOptions []string `json:"questionOptions"`
}
