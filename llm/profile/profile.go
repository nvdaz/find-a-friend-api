package profile

import (
	"encoding/json"

	"github.com/nvdaz/find-a-friend-api/model"
)

func GenerateProfile(questions []string) (*model.InternalProfile, error) {
	data, err := json.Marshal(questions)
	if err != nil {
		return nil, err
	}

	intermediateProfile, err := initializeProfile(string(data))
	if err != nil {
		return nil, err
	}

	// err = reviseProfile(intermediateProfile)
	// if err != nil {
	// 	return nil, err
	// }

	nonSecretProfile := model.NewNonSecretIntermediateProfile(intermediateProfile)

	features, err := generateUserFeatures(nonSecretProfile)
	if err != nil {
		return nil, err
	}

	return &model.InternalProfile{
		Interests:                intermediateProfile.Interests,
		Personality:              intermediateProfile.Personality,
		Skills:                   intermediateProfile.Skills,
		Goals:                    intermediateProfile.Goals,
		Values:                   intermediateProfile.Values,
		Demographics:             intermediateProfile.Demographics,
		LivedExperiences:         intermediateProfile.LivedExperiences,
		Habits:                   intermediateProfile.Habits,
		InterpersonalSkills:      intermediateProfile.InterpersonalSkills,
		ExceptionalCircumstances: intermediateProfile.ExceptionalCircumstances,
		Summary:                  features.Summary,
		Tags:                     features.Tags,
		Bio:                      features.Bio,
		KeyQuestions:             features.KeyQuestions,
	}, nil

}
