package profile

import (
	"encoding/json"
	"fmt"

	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/model"
)

func GenerateProfile(id string, questions []string, conversations [][]db.Message) (*model.InternalProfile, error) {
	data, err := json.Marshal(questions)
	if err != nil {
		return nil, err
	}

	simplifiedConversations := make([][]string, 0, len(conversations))
	for _, conversation := range conversations {
		simplifiedConversation := make([]string, 0, len(conversation))
		for _, message := range conversation {
			simplifiedConversation = append(simplifiedConversation, fmt.Sprintf("%s: %s", message.SenderId, message.Message))
		}
		simplifiedConversations = append(simplifiedConversations, simplifiedConversation)
	}

	conversationData, err := json.Marshal(simplifiedConversations)
	if err != nil {
		return nil, err
	}

	intermediateProfile, err := initializeProfile(id, string(data), string(conversationData))
	if err != nil {
		return nil, err
	}

	// err = reviseProfile(intermediateProfile)
	// if err != nil {
	// 	return nil, err
	// }

	features, err := generateUserFeatures(*intermediateProfile, string(data))
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
		Hobbies:                  intermediateProfile.Hobbies,
		InterpersonalSkills:      intermediateProfile.InterpersonalSkills,
		ExceptionalCircumstances: intermediateProfile.ExceptionalCircumstances,
		Topics:                   intermediateProfile.Topics,
		Summary:                  features.Summary,
		Tags:                     features.Tags,
		Bio:                      features.Bio,
		KeyQuestions:             features.KeyQuestions,
		Subtitle:                 features.Subtitle,
		LookingFor:               features.LookingFor,
	}, nil

}
