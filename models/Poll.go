package models

//Poll This object contains information about a poll.
type Poll struct {
	ID                    string           `json:"id"`
	Question              string           `json:"question"`
	Options               []*PollOption    `json:"options"`
	TotalVoterCounter     int              `json:"total_voter_count"`
	IsClosed              bool             `json:"is_closed"`
	IsAnonymous           bool             `json:"is_anonymous"`
	Type                  string           `json:"type"`
	AllowsMultipleAnswers bool             `json:"allows_multiple_answers"`
	CorrectOptionID       int              `json:"correct_option_id"`
	Explanation           string           `json:"explanation"`
	ExplanationEntities   []*MessageEntity `json:"explanation_entities"`
	OpenPeriod            int              `json:"open_period"`
	CloseDate             int              `json:"close_date"`
}
