package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/llm/profile"
	"github.com/nvdaz/find-a-friend-api/model"
)

type UserService struct {
	userStore                db.UserStore
	serviceConversationStore db.ServiceConversationStore
}

func NewUserService(userStore db.UserStore, serviceConversationStore db.ServiceConversationStore) UserService {
	return UserService{userStore, serviceConversationStore}
}

func needsUpdate(user *db.User) bool {
	updated, err := time.Parse(time.RFC3339, user.UpdatedAt)
	if err != nil {
		return true
	}

	if user.GeneratedAt == nil || user.Profile == nil {
		return true

	}

	refreshed, err := time.Parse(time.RFC3339, *user.GeneratedAt)
	if err != nil {
		return true
	}

	if updated.Sub(refreshed) < time.Minute {
		return false
	}

	return true

}

func (service *UserService) MarkUserAsUpdated(id string) error {
	return service.userStore.MarkUserAsUpdated(id)
}

func (service *UserService) GetUser(id string) (*model.User, error) {
	user, err := service.userStore.GetUser(id)
	if err != nil {
		return nil, err
	}

	if !needsUpdate(user) {
		profile := model.InternalProfile{}
		err = json.Unmarshal([]byte(*user.Profile), &profile)
		if err != nil {
			return nil, err
		}

		fmt.Println("Profile", profile)

		return &model.User{
			Id:      user.Id,
			Name:    user.Name,
			Profile: &profile,
		}, nil
	}

	serviceConversations, err := service.serviceConversationStore.GetRecentServiceConversations(id, 100)
	if err != nil {
		return nil, err
	}

	questions := make([]string, 0)
	for _, conversation := range serviceConversations {
		questions = append(questions, conversation.Question)
	}

	profile, err := profile.GenerateProfile(questions)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(profile)
	if err != nil {
		return nil, err
	}

	service.userStore.UpdateUserProfile(id, string(data))

	return &model.User{
		Id:      user.Id,
		Name:    user.Name,
		Profile: profile,
	}, nil
}

func (service *UserService) GetAllUsers() ([]model.User, error) {
	users, err := service.userStore.GetAllUsers()
	if err != nil {
		return nil, err
	}

	result := make([]model.User, 0)
	for _, user := range users {
		if user.Profile == nil {
			continue
		}

		profile := model.InternalProfile{}
		err = json.Unmarshal([]byte(*user.Profile), &profile)
		if err != nil {
			return nil, err
		}

		result = append(result, model.User{
			Id:      user.Id,
			Name:    user.Name,
			Profile: &profile,
		})
	}

	return result, nil
}

func (service *UserService) CreateUser(createUser db.CreateUser) error {
	return service.userStore.CreateUser(createUser)
}
