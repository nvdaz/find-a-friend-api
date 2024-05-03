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
	sent, err := service.messagesStore.GetRecentMessages(senderId, receiverId, limit)
	if err != nil {
		return nil, err
	}
	received, err := service.messagesStore.GetRecentMessages(receiverId, senderId, limit)
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

func (service *MessageService) PollMessages(senderId, receiverId, after string, limit int) ([]model.Message, error) {
	sent, err := service.messagesStore.GetNewMessages(senderId, receiverId, after, limit)
	if err != nil {
		return nil, err
	}
	received, err := service.messagesStore.GetNewMessages(receiverId, senderId, after, limit)
	if err != nil {
		return nil, err
	}

	conversation := make([]model.Message, 0)
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
