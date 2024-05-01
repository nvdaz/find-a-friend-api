package match

import (
	"encoding/json"
	"fmt"

	"github.com/nvdaz/find-a-friend-api/llm"
	"github.com/nvdaz/find-a-friend-api/model"
)

func ExplainMatch(user1, user2 model.User) (string, error) {
	data := struct {
		User1 model.User `json:"user1"`
		User2 model.User `json:"user2"`
	}{
		User1: user1,
		User2: user2,
	}

	prompt, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	explanation := struct {
		Explanation string `json:"explanation"`
	}{}

	err = llm.GetResponseJson(&explanation, llm.ModelClaudeSonnet, string(prompt), "Your job is to explain why these two users are a good match. Go into as much detail as possible with a 200 word justifications. Respond with a JSON object without formatting containing a single key 'explanation', which is a string that explains why these two users are a good match.", nil)

	return explanation.Explanation, err
}

func DecideBestMatch(explanations map[string]string) (string, error) {
	prompt, err := json.Marshal(explanations)
	if err != nil {
		return "", err
	}

	bestMatch := struct {
		BestMatch string `json:"best_match"`
	}{}

	err = llm.GetResponseJson(&bestMatch, llm.ModelGpt4, string(prompt), "Your job is to decide which of the potential matches is the best match based on the explanations provided. Respond with a JSON object without formatting containing a single key 'best_match', which is the ID of the best match.", nil)

	return bestMatch.BestMatch, err
}

func ExplainMatchToUser(user1, user2 model.User) (string, error) {
	data := struct {
		User1 model.User `json:"user1"`
		User2 model.User `json:"user2"`
	}{
		User1: user1,
		User2: user2,
	}

	prompt, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	explanation := struct {
		Explanation string `json:"explanation"`
	}{}

	err = llm.GetResponseJson(&explanation, llm.ModelGpt3p5, string(prompt), fmt.Sprintf("You are a matchmaker. Write a personalized message to %q (refer to them as 'you') why %q would be a good friend for them. Go into as much detail as possible with a 1-paragraph, 60 word justification. Be sure to use the matched user's name and specific details about their profile in your explanation. Use casual, friendly language. Respond with a JSON object without formatting containing a single key 'explanation'.", user1.Name, user2.Name), nil)
	if err != nil {
		return "", err
	}

	return explanation.Explanation, nil
}
