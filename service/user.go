package service

import (
	"fmt"
	"log"
	"time"

	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/llm"
	"github.com/nvdaz/find-a-friend-api/model"
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

func buildUser(user *db.User, userProfile *db.UserProfile, interests []db.Interest, keyQuestions []db.ServiceConversation) model.User {
	convertedInterests := []model.Interest{}
	for _, interest := range interests {
		convertedInterests = append(convertedInterests, model.Interest(interest))
	}

	convertedKeyQuestions := []string{}
	for _, conversation := range keyQuestions {
		convertedKeyQuestions = append(convertedKeyQuestions, conversation.Question)
	}

	modelUser := model.User{
		Id:           user.Id,
		Name:         user.Name,
		Bio:          userProfile.Bio,
		Personality:  model.Personality(userProfile.Personality),
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

func (service *UserService) GetUser(id string) (*model.User, error) {
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

	generalInterests, err := llm.ExtrapolateUserInterests(interests)
	if err != nil {
		return nil, err
	}
	log.Println("General Interests", generalInterests)

	aggregateInterests := append(interests, generalInterests...)
	service.interestsStore.InsertUserInterests(id, aggregateInterests)

	personality, err := llm.GenerateUserPersonality(serviceConversations)
	if err != nil {
		return nil, err
	}
	userProfile.Personality = *personality
	fmt.Println("Personality", personality)

	bio, err := llm.GenerateUserBio(*user, *personality, aggregateInterests)
	if err != nil {
		return nil, err
	}
	userProfile.Bio = *bio
	fmt.Println("Bio", bio)

	keyConversations, err := llm.GenerateKeyServiceConversations(*user, serviceConversations)
	if err != nil {
		return nil, err
	}

	err = service.serviceConversationStore.UpdateKeyServiceConversations(id, keyConversations)
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

	modelUser := buildUser(user, userProfile, aggregateInterests, keyServiceConversations)

	return &modelUser, nil
}

func (service *UserService) GetAllUsers() ([]model.User, error) {
	userProfiles, err := service.userProfileStore.GetAllUserProfiles()
	if err != nil {
		return nil, err
	}

	users := []model.User{}

	for _, userProfile := range userProfiles {
		user, err := service.userStore.GetUser(userProfile.Id)
		if err != nil {
			return nil, err
		}

		keyServiceConversations, err := service.serviceConversationStore.GetKeyServiceConversations(userProfile.Id)
		if err != nil {
			return nil, err
		}

		interests, err := service.interestsStore.GetUserInterests(userProfile.Id)
		if err != nil {
			return nil, err
		}

		modelUser := buildUser(user, &userProfile, interests, keyServiceConversations)

		users = append(users, modelUser)
	}

	return users, nil
}

func (service *UserService) GetBestMatch(userId string) (*model.Match, error) {
	allUsers, err := service.GetAllUsers()
	if err != nil {
		return nil, err
	}

	user := model.User{}
	users := []model.User{}
	for _, tUser := range allUsers {
		if tUser.Id != userId {
			users = append(users, tUser)
		} else {
			user = tUser
		}
	}

	if user.Id == "" {
		return nil, nil
	}

	match, err := llm.GenerateUserBestMatch(user, users)
	if err != nil {
		return nil, err
	}

	return match, nil
}
