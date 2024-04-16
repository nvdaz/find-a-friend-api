package llm

import (
	"encoding/json"
	"fmt"

	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/model"
)

func GenerateUserBio(user db.User, personality db.Personality, interests []db.Interest) (*string, error) {
	system := "Create a short introductory biography from the perspective of the provided user using personal pronouns. Include a brief description of their personality and interests. The biography should be a single paragraph, no more than 200 words in length. Provide a JSON object without any formatting containing two keys: 'bio', with the value being the biography."

	prompt, err := json.Marshal(struct {
		Name        string         `json:"name"`
		Personality db.Personality `json:"personality"`
		Interests   []db.Interest  `json:"interests"`
	}{user.Name, personality, interests})
	if err != nil {
		return nil, err
	}

	result := struct {
		Bio string `json:"bio"`
	}{}
	if err = GetResponseJson(&result, ModelGpt4, string(prompt), system, nil); err != nil {
		return nil, err
	}

	return &result.Bio, nil
}

const UserInterestsCount = 15

func GenerateUserInterests(serviceConversations []db.ServiceConversation) ([]db.Interest, error) {
	system := fmt.Sprintf("Create a list of interests based on the provided service conversations. The list should contain at least %d specific interests. Provide a JSON object without any formatting containing the keys 'interests', with the value being the list of interests. The interests should be objects with a key 'interest' containing the interest, a key 'skill' containing the skill level on a scale of 1 to 5, and a key 'intensity' containing the intensity of the interest on a scale of 1 to 5.", UserInterestsCount)
	questions := make([]string, 0)
	for _, conversation := range serviceConversations {
		questions = append(questions, conversation.Question)
	}
	prompt, err := json.Marshal(questions)
	if err != nil {
		return nil, err
	}

	result := struct {
		Interests []db.Interest `json:"interests"`
	}{}
	if err = GetResponseJson(&result, ModelGpt4, string(prompt), system, nil); err != nil {
		return nil, err
	}

	return result.Interests, nil
}

const GeneralUserInterestsCount = 5

func ExtrapolateUserInterests(interests []db.Interest) ([]db.Interest, error) {
	system := fmt.Sprintf("Create a list of general interests based on the provided list of specific interests. The list should contain at least %d general interests. Provide a JSON object without any formatting containing the keys 'interests', with the value being the list of interests. The interests should be objects with a key 'interest' containing the interest, a key 'skill' containing the skill level on a scale of 1 to 5, and a key 'intensity' containing the intensity of the interest on a scale of 1 to 5.", GeneralUserInterestsCount)
	prompt, err := json.Marshal(interests)
	if err != nil {
		return nil, err
	}

	result := struct {
		Interests []db.Interest `json:"interests"`
	}{}
	if err = GetResponseJson(&result, ModelGpt4, string(prompt), system, nil); err != nil {
		return nil, err
	}

	return result.Interests, nil
}

func GenerateUserPersonality(serviceConversations []db.ServiceConversation) (*db.Personality, error) {
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
	if err = GetResponseJson(&result, ModelGpt4, string(prompt), system, nil); err != nil {
		return nil, err
	}

	return &result, nil
}

func GenerateKeyServiceConversations(user db.User, serviceConversations []db.ServiceConversation) ([]string, error) {
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

	result := struct {
		Conversations []string `json:"key_conversations"`
	}{}
	if err = GetResponseJson(&result, ModelGpt4, string(prompt), system, nil); err != nil {
		return nil, err
	}

	return result.Conversations, nil
}

func GenerateUserBestMatch(user model.User, otherUsers []model.User) (*model.Match, error) {
	system := "Determine the user that the provided user is most likely to be friends with based on their shared personalities and interests. You are provided with the user object as well as a list of other users. Provide ONLY a JSON object with keys 'user_id' with the id of the matched user, 'affinity' with a numeric affinity score, and 'reason' with a brief explanation of why the match was made."
	type GenerateUserAffinitiesPrompt struct {
		User       model.User   `json:"user"`
		OtherUsers []model.User `json:"other_users"`
	}
	prompt, err := json.Marshal(GenerateUserAffinitiesPrompt{user, otherUsers})
	if err != nil {
		return nil, err
	}

	result := model.Match{}
	if err = GetResponseJson(&result, ModelGpt4, string(prompt), system, nil); err != nil {
		return nil, err
	}

	return &result, nil
}
