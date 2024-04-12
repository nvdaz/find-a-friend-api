package llm

import (
	"encoding/json"

	"github.com/nvdaz/find-a-friend-api/db"
)

type GenerateUserBioResult struct {
	Bio string `json:"bio"`
}

func GenerateUserBio(user db.User) (*string, error) {
	model := "gpt3-5"
	system := "Create a short biography for the provided user. Include a brief description of their personality and interests. The biography should be a single paragraph, no more than 200 words in length. Provide a JSON object without any formatting containing two keys: 'bio', with the value being the biography."
	prompt, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	result := GenerateUserBioResult{}
	if err = GetResponseJson(&result, model, string(prompt), system, nil); err != nil {
		return nil, err
	}

	return &result.Bio, nil
}
