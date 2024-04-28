package model

type Match struct {
	Id      string `json:"id"`
	UserId  string `json:"user_id"`
	MatchId string `json:"other_id"`
	Reason  string `json:"reason"`
}
