package profile

import (
	"fmt"

	"github.com/nvdaz/find-a-friend-api/llm"
	"github.com/nvdaz/find-a-friend-api/model"
)

func GeneratePersonalityFromConversations(id, conversations string) (model.Personality, error) {
	system := fmt.Sprintf("Let us play a guessing game. You are provided with a list of conversations the user had with other users. Your task is to guess the user's (%s) personality based on the Big Five (OCEAN) model, assigning scores from 0 to 5 for each trait, where 0 means the trait is not present and 5 signifies a strong presence. List the scores for Openness, Conscientiousness, Extroersion, Agreeableness, and Neuroticism. Start with analysis and use deductive reasoning to answer as precisely as possible. Then, provide a JSON object in a JSON code block containing the keys 'openness', 'conscientiousness', 'extroversion', 'agreeableness', and 'neuroticism'.", id)

	personality := model.Personality{}

	err := llm.GetResponseJson(&personality, llm.ModelClaudeSonnet, conversations, system, nil)
	if err != nil {
		return model.Personality{}, nil
	}

	return personality, nil
}

func GenerateInterpersonalSkillsFromConversations(id, conversations string) (model.InterpersonalSkills, error) {
	system := fmt.Sprintf("Let us play a guessing game. You are provided with a list of conversations the user had with other users. Your task is to guess the user's (%s) interpersonal skills based on their conversations with others. Start with analysis and use deductive reasoning to answer as precisely as possible. Then, provide a JSON object in a JSON code block containing the keys 'active_listening', 'teamwork', 'responsibility', 'dependability', 'leadership', 'motivation', 'flexibility', 'patience', and 'empathy'. Each key should have a value between 0 and 1, representing the strength of the skill.", id)

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
