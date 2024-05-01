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

func (service *MessageService) GetMessages(senderId, receiverId string) ([]model.Message, error) {
	sent, err := service.messagesStore.GetRecentMessages(senderId, receiverId, 50)
	if err != nil {
		return nil, err
	}
	received, err := service.messagesStore.GetRecentMessages(receiverId, senderId, 50)
	if err != nil {
		return nil, err
	}

	conversation := make([]model.Message, 0, 100)
	sentIndex := 0
	receivedIndex := 0
	for sentIndex < len(sent) && receivedIndex < len(received) {
		if sent[sentIndex].CreatedAt > received[receivedIndex].CreatedAt {
			conversation = append(conversation, model.Message(sent[sentIndex]))
			sentIndex++
		} else {
			conversation = append(conversation, model.Message(received[receivedIndex]))
			receivedIndex++
		}
	}

	for sentIndex < len(sent) {
		conversation = append(conversation, model.Message(sent[sentIndex]))
		sentIndex++
	}

	for receivedIndex < len(received) {
		conversation = append(conversation, model.Message(received[receivedIndex]))
		receivedIndex++
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
