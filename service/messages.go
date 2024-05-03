package service

import (
	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/model"
)

type MessageService struct {
	messagesStore db.MessageStore
	userService   UserService
}

func NewMessagesService(messagesStore db.MessageStore, userService UserService) MessageService {
	return MessageService{messagesStore, userService}
}

func (service *MessageService) GetMessages(senderId, receiverId string, limit int) ([]model.Message, error) {
	messages, err := service.messagesStore.GetRecentMessages(senderId, receiverId, limit)
	if err != nil {
		return nil, err
	}

	conversation := make([]model.Message, 0, len(messages))
	for _, message := range messages {
		conversation = append(conversation, model.Message(message))
	}

	return conversation, nil
}

func (service *MessageService) PollMessages(senderId, receiverId, after string, limit int) ([]model.Message, error) {
	messages, err := service.messagesStore.GetNewMessages(senderId, receiverId, after, limit)
	if err != nil {
		return nil, err
	}

	conversation := make([]model.Message, 0, len(messages))
	for _, message := range messages {
		conversation = append(conversation, model.Message(message))
	}

	return conversation, nil
}

func (service *MessageService) CreateMessage(senderId, receiverId, message string) error {
	err := service.messagesStore.CreateMessage(senderId, receiverId, message)
	if err != nil {
		return err
	}

	service.userService.MarkUserAsUpdated(senderId)
	service.userService.MarkUserAsUpdated(receiverId)

	return nil
}


