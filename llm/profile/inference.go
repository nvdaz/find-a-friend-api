package profile

import (
	"encoding/json"

	"github.com/nvdaz/find-a-friend-api/llm"
	"github.com/nvdaz/find-a-friend-api/model"
)

func reviseProfile(user *model.IntermediateProfile) error {
	profileString, err := json.Marshal(user)
	if err != nil {
		return nil
	}

	err = llm.GetResponseJson(&user, llm.ModelGpt3p5, string(profileString), "Your job is to revise the profile, inferring any missing information. Respond with the updated profile in JSON format exactly in the format it was received.", nil)
	if err != nil {
		return err
	}

	return nil
}
