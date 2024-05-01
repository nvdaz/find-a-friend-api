package match

import (
	"encoding/json"
	"fmt"

	"github.com/nvdaz/find-a-friend-api/llm"
	"github.com/nvdaz/find-a-friend-api/model"
)

const CandidateMatchesCount = 4

func GenerateCandidateMatches(user model.User, users []model.User) ([]string, error) {
	type UserSummary struct {
		Id      string `json:"id"`
		Summary string `json:"summary"`
	}

	var userSummaries []UserSummary
	for _, user := range users {
		userSummaries = append(userSummaries, UserSummary{Id: user.Id, Summary: user.Profile.Summary})
	}

	d := struct {
		User  model.User    `json:"user"`
		Users []UserSummary `json:"users"`
	}{
		User:  user,
		Users: userSummaries,
	}

	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	matches := struct {
		Matches []string `json:"matches"`
	}{}

	err = llm.GetResponseJson(&matches, llm.ModelClaudeSonnet, string(data), fmt.Sprintf("Your job is to generate a list of %d potential matches based on the user summaries provided. Respond with JSON with a key 'matches', a list of user IDs that are potential matches.", CandidateMatchesCount), nil)
	if err != nil {
		return nil, err
	}

	return matches.Matches, nil
}
