package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/llm/profile"
	"github.com/nvdaz/find-a-friend-api/model"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userStore    db.UserStore
	messageStore db.MessageStore
}

func NewUserService(userStore db.UserStore, messageStore db.MessageStore) UserService {
	return UserService{userStore, messageStore}
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
			Avatar:  user.Avatar,
			Profile: &profile,
		}, nil
	}

	questions, err := service.GetAgentQuestions(id, 60)
	if err != nil {
		return nil, err
	}

	conversations, err := service.messageStore.GetRecentMessagesAllConversations(id, 60)
	if err != nil {
		return nil, err
	}
	partitionedConversations := partitionConversations(id, conversations)

	profile, err := profile.GenerateProfile(id, questions, partitionedConversations)
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
		Avatar:  user.Avatar,
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

type RegisterUser struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (service *UserService) RegisterUser(registerUser RegisterUser) (*model.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(registerUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()

	err = service.userStore.CreateUser(db.CreateUser{
		Id:       id,
		Name:     registerUser.Name,
		Username: registerUser.Username,
		Password: string(password),
	})
	if err != nil {
		return nil, err
	}

	return &model.User{
		Id:   id,
		Name: registerUser.Name,
	}, nil

}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (service *UserService) LoginUser(loginUser LoginUser) (*model.User, error) {
	user, err := service.userStore.GetUserByUsername(loginUser.Username)
	if err != nil {
		fmt.Println("Error getting user by username", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		return nil, err
	}

	var profile *model.InternalProfile
	if user.Profile != nil {
		profile = &model.InternalProfile{}
		json.Unmarshal([]byte(*user.Profile), profile)
	}

	fmt.Println("Profile", *user.Profile)

	return &model.User{
		Id:      user.Id,
		Name:    user.Name,
		Avatar:  user.Avatar,
		Profile: profile,
	}, nil
}

func (service *UserService) UpdateAvatar(id, avatar string) error {
	return service.userStore.UpdateAvatar(id, avatar)
}

func (service *UserService) GetAgentQuestions(id string, limit int) ([]string, error) {
	questions, err := service.messageStore.GetRecentSentMessages(id, "00000000-0000-0000-0000-000000000000", limit)
	if err != nil {
		return nil, err
	}

	questionStrings := make([]string, len(questions))
	for i, question := range questions {
		questionStrings[i] = question.Message
	}

	return questionStrings, nil
}

func partitionConversations(id string, messages []db.Message) [][]db.Message {
	conversations := map[string][]db.Message{}

	for _, message := range messages {
		conversationId := message.SenderId
		if message.SenderId == id {
			conversationId = message.ReceiverId
		}

		if _, ok := conversations[conversationId]; !ok {
			conversations[conversationId] = []db.Message{}
		}
		conversations[conversationId] = append(conversations[conversationId], message)
	}

	partitions := [][]db.Message{}
	for _, messages := range conversations {
		partitions = append(partitions, messages)
	}

	return partitions
}
