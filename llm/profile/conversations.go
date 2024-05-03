package profile

import (
	"fmt"

	"github.com/nvdaz/find-a-friend-api/llm"
	"github.com/nvdaz/find-a-friend-api/model"
)

func GeneratePersonalityFromConversations(id, conversations string) (model.Personality, error) {
	system := fmt.Sprintf("Assess the user's (%s) preliminary personality based on the Big Five (OCEAN) model, assigning scores from 0 to 5 for each trait, where 0 means the trait is not present and 5 signifies a strong presence. You are provided with a list of conversations the user had with other users. If it is not possible to determine a trait, provide an average value. List the scores for Openness, Conscientiousness, Extroersion, Agreeableness, and Neuroticism. These scores are only preliminary and need not be perfectly accurate. Provide a JSON object without any formatting containing the keys 'openness', 'conscientiousness', 'extroversion', 'agreeableness', and 'neuroticism'. Start with an in-depth analysis of the user's queries in an 'analysis' key. ", id)

	personality := model.Personality{}

	err := llm.GetResponseJson(&personality, llm.ModelClaudeSonnet, conversations, system, nil)
	if err != nil {
		return model.Personality{}, nil
	}

	return personality, nil
}

func GenerateInterpersonalSkillsFromConversations(id, conversations string) (model.InterpersonalSkills, error) {
	system := fmt.Sprintf("Assess the user's (%s) preliminary interpersonal skills based on their conversations with others. You are provided with a list of conversations the user had with other users. Provide a JSON object without any formatting containing the keys 'active_listening', 'teamwork', 'responsibility', 'dependability', 'leadership', 'motivation', 'flexibility', 'patience', and 'empathy'. Each key should have a value between 0 and 1, representing the strength of the skill. Start with an in-depth analysis of the user's queries in an 'analysis' key. ", id)

	interpersonalSkills := model.InterpersonalSkills{}

	err := llm.GetResponseJson(&interpersonalSkills, llm.ModelClaudeSonnet, conversations, system, nil)
	if err != nil {
		return model.InterpersonalSkills{}, nil
	}

	return interpersonalSkills, nil
}

func GenerateTopicsFromConversations(conversations string) ([]model.Topic, error) {
	system := "Summarize the topics of conversation based on the user's conversations with others. You are provided with a list of conversations the user had with other users. Provide a JSON object without any formatting containing a key 'topics', with the value being a list of topics discussed in the conversations. Each topic should have a 'topic' key with the topic name, a 'level' key with a value between 0 and 1 representing the importance of the topic, and an 'emoji' key with an emoji representing the topic."

	topics := struct {
		Topics []model.Topic `json:"topics"`
	}{}

	err := llm.GetResponseJson(&topics, llm.ModelClaudeSonnet, conversations, system, nil)
	if err != nil {
		return []model.Topic{}, nil
	}

	return topics.Topics, nil
}
