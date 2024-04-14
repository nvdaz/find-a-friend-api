package service

import (
	"fmt"
	"time"

	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/llm"
)

type UserService struct {
	userStore                db.UserStore
	serviceConversationStore db.ServiceConversationStore
	interestsStore           db.InterestsStore
	userProfileStore         db.UserProfilesStore
}

func NewUserService(userStore db.UserStore, serviceConversationStore db.ServiceConversationStore, interestsStore db.InterestsStore, userDerivedStore db.UserProfilesStore) UserService {
	return UserService{userStore, serviceConversationStore, interestsStore, userDerivedStore}
}

type Personality struct {
	Extroversion      float64 `json:"extroversion"`
	Agreeableness     float64 `json:"agreeableness"`
	Conscientiousness float64 `json:"conscientiousness"`
	Neuroticism       float64 `json:"neuroticism"`
	Openness          float64 `json:"openness"`
}

type Interest struct {
	Interest  string  `json:"interest"`
	Intensity float64 `json:"intensity"`
	Skill     float64 `json:"skill"`
}

type User struct {
	Id           string      `json:"id"`
	Name         string      `json:"name"`
	Bio          string      `json:"bio"`
	Personality  Personality `json:"personality"`
	Interests    []Interest  `json:"interests"`
	KeyQuestions []string    `json:"key_questions"`
}

func buildUser(user *db.User, userProfile *db.UserProfile, interests []db.Interest, keyQuestions []db.ServiceConversation) User {
	convertedInterests := []Interest{}
	for _, interest := range interests {
		convertedInterests = append(convertedInterests, Interest(interest))
	}

	convertedKeyQuestions := []string{}
	for _, conversation := range keyQuestions {
		convertedKeyQuestions = append(convertedKeyQuestions, conversation.Question)
	}

	modelUser := User{
		Id:           user.Id,
		Name:         user.Name,
		Bio:          userProfile.Bio,
		Personality:  Personality(userProfile.Personality),
		Interests:    convertedInterests,
		KeyQuestions: convertedKeyQuestions,
	}

	return modelUser
}

func needsUpdate(user *db.User, userProfile *db.UserProfile) bool {
	updated, err := time.Parse(time.RFC3339, user.UpdatedAt)
	if err != nil {
		return true
	}

	refreshed, err := time.Parse(time.RFC3339, userProfile.UpdatedAt)
	if err != nil {
		return true
	}

	if updated.Sub(refreshed) < time.Minute {
		return false
	}

	return true

}

func (service *UserService) GetUser(id string) (*User, error) {
	user, err := service.userStore.GetUser(id)
	if err != nil {
		return nil, err
	}
	userProfile, err := service.userProfileStore.GetUserProfile(id)
	if err != nil {
		if err != db.ErrUserProfileNotFound {
			return nil, err
		}

		err = nil
		userProfile = &db.UserProfile{
			Id: id,
		}
	}

	if !needsUpdate(user, userProfile) {
		keyServiceConversations, err := service.serviceConversationStore.GetKeyServiceConversations(id)
		if err != nil {
			return nil, err
		}

		interests, err := service.interestsStore.GetUserInterests(id)
		if err != nil {
			return nil, err
		}

		modelUser := buildUser(user, userProfile, interests, keyServiceConversations)

		return &modelUser, nil
	}

	serviceConversations, err := service.serviceConversationStore.GetRecentServiceConversations(id, 100)
	if err != nil {
		return nil, err
	}

	interests, err := llm.GenerateUserInterests(serviceConversations)
	if err != nil {
		return nil, err
	}
	service.interestsStore.InsertUserInterests(id, interests)
	fmt.Println("Interests", interests)

	personality, err := llm.GenerateUserPersonality(serviceConversations)
	if err != nil {
		return nil, err
	}
	userProfile.Personality = *personality
	fmt.Println("Personality", personality)

	bio, err := llm.GenerateUserBio(*user, *personality, interests)
	if err != nil {
		return nil, err
	}
	userProfile.Bio = *bio
	fmt.Println("Bio", bio)

	key_conversations, err := llm.GenerateKeyServiceConversations(*user, serviceConversations)
	if err != nil {
		return nil, err
	}

	err = service.serviceConversationStore.UpdateKeyServiceConversations(id, key_conversations)
	if err != nil {
		return nil, err
	}

	err = service.userProfileStore.InsertUserProfile(*userProfile)
	if err != nil {
		return nil, err
	}

	keyServiceConversations, err := service.serviceConversationStore.GetKeyServiceConversations(id)
	if err != nil {
		return nil, err
	}
	fmt.Println("Key Conversations", keyServiceConversations)

	modelUser := buildUser(user, userProfile, interests, keyServiceConversations)

	return &modelUser, nil
}
