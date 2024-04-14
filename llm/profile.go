package llm

import (
	"encoding/json"

	"github.com/nvdaz/find-a-friend-api/db"
)

func GenerateUserBio(user db.User, personality db.Personality, interests []db.Interest) (*string, error) {
	model := "gpt4-new"
	system := "Create a short introductory biography from the perspective of the provided user using personal pronouns. Include a brief description of their personality and interests. The biography should be a single paragraph, no more than 200 words in length. Provide a JSON object without any formatting containing two keys: 'bio', with the value being the biography."

	type GenerateUserBioPrompt struct {
		Name        string         `json:"name"`
		Personality db.Personality `json:"personality"`
		Interests   []db.Interest  `json:"interests"`
	}

	prompt, err := json.Marshal(GenerateUserBioPrompt{user.Name, personality, interests})
	if err != nil {
		return nil, err
	}

	type GenerateUserBioResult struct {
		Bio string `json:"bio"`
	}
	result := GenerateUserBioResult{}
	if err = GetResponseJson(&result, model, string(prompt), system, nil); err != nil {
		return nil, err
	}

	return &result.Bio, nil
}

func GenerateUserInterests(serviceConversations []db.ServiceConversation) ([]db.Interest, error) {
	model := "gpt4-new"
	system := "Create a list of interests based on the provided service conversations. The list should contain at least 15 specific interests. Provide a JSON object without any formatting containing the keys 'interests', with the value being the list of interests. The interests should be objects with a key 'interest' containing the interest, a key 'skill' containing the skill level on a scale of 1 to 5, and a key 'intensity' containing the intensity of the interest on a scale of 1 to 5."
	questions := make([]string, 0)
	for _, conversation := range serviceConversations {
		questions = append(questions, conversation.Question)
	}
	prompt, err := json.Marshal(questions)
	if err != nil {
		return nil, err
	}

	type GenerateUserInterestsResult struct {
		Interests []db.Interest `json:"interests"`
	}
	result := GenerateUserInterestsResult{}
	if err = GetResponseJson(&result, model, string(prompt), system, nil); err != nil {
		return nil, err
	}

	return result.Interests, nil
}

func GenerateUserPersonality(serviceConversations []db.ServiceConversation) (*db.Personality, error) {
	model := "gpt4-new"
	system := "Assess the user's preliminary personality based on the Big Five (OCEAN) model, assigning scores from 0 to 5 for each trait, where 0 means the trait is not present and 5 signifies a strong presence. If it is not possible to determine a trait, provide an average value. List the scores for Openness, Conscientiousness, Extraversion, Agreeableness, and Neuroticism. These scores are only preliminary and need not be perfectly accurate. Provide a JSON object without any formatting containing the keys 'openness', 'conscientiousness', 'extraversion', 'agreeableness', and 'neuroticism'."
	questions := make([]string, 0)
	for _, conversation := range serviceConversations {
		questions = append(questions, conversation.Question)
	}
	prompt, err := json.Marshal(questions)
	if err != nil {
		return nil, err
	}

	result := db.Personality{}
	if err = GetResponseJson(&result, model, string(prompt), system, nil); err != nil {
		return nil, err
	}

	return &result, nil
}

func GenerateKeyServiceConversations(user db.User, serviceConversations []db.ServiceConversation) ([]string, error) {
	model := "gpt4-new"
	system := "Identify three key service conversations from the provided list of service conversations. Key service conversations are those that are most important or relevant to the user. You are provided with the user object as well as a list of questions they asked. Provide a JSON object without any formatting containing the key 'key_conversations', with the value being a list of the ids of key service conversations."
	type GenerateKeyServiceConversationsPromptQuestion struct {
		Id       string `json:"id"`
		Question string `json:"question"`
	}
	questions := make([]GenerateKeyServiceConversationsPromptQuestion, 0)
	for _, conversation := range serviceConversations {
		questions = append(questions, GenerateKeyServiceConversationsPromptQuestion{conversation.Id, conversation.Question})
	}
	type GenerateKeyServiceConversationsPrompt struct {
		User      db.User                                         `json:"user"`
		Questions []GenerateKeyServiceConversationsPromptQuestion `json:"questions"`
	}
	info := GenerateKeyServiceConversationsPrompt{user, questions}
	prompt, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	type GenerateKeyServiceConversationsResult struct {
		Conversations []string `json:"key_conversations"`
	}
	result := GenerateKeyServiceConversationsResult{}
	if err = GetResponseJson(&result, model, string(prompt), system, nil); err != nil {
		return nil, err
	}

	return result.Conversations, nil
}
