package db

import (
	"database/sql"

	"github.com/google/uuid"
)

type MessageStore struct {
	db *sql.DB
}

func NewMessagesStore(db *sql.DB) MessageStore {
	return MessageStore{db}
}

type Message struct {
	Id         string
	SenderId   string
	ReceiverId string
	Message    string
	CreatedAt  string
}

func (store *MessageStore) GetRecentMessages(senderId, receiverId string, limit int) ([]Message, error) {
	rows, err := store.db.Query(
		`SELECT id, sender_id, receiver_id, message, created_at
		 FROM messages
		 WHERE sender_id = ? AND receiver_id = ?
		 ORDER BY created_at DESC
		 LIMIT ?`,
		senderId, receiverId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userConversations := []Message{}
	for rows.Next() {
		userConversation := Message{}
		if err := rows.Scan(&userConversation.Id, &userConversation.SenderId,
			&userConversation.ReceiverId, &userConversation.Message, &userConversation.CreatedAt); err != nil {
			return nil, err
		}
		userConversations = append(userConversations, userConversation)
	}

	return userConversations, nil
}

func (store *MessageStore) GetNewMessages(senderId, receiverId string, after string, limit int) ([]Message, error) {
	rows, err := store.db.Query(
		`SELECT id, sender_id, receiver_id, message, created_at
		 FROM messages
		 WHERE sender_id = ? AND receiver_id = ? AND created_at > datetime(?)
		 ORDER BY created_at DESC
		 LIMIT ?
		 `,
		senderId, receiverId, after, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userConversations := []Message{}
	for rows.Next() {
		userConversation := Message{}
		if err := rows.Scan(&userConversation.Id, &userConversation.SenderId,
			&userConversation.ReceiverId, &userConversation.Message, &userConversation.CreatedAt); err != nil {
			return nil, err
		}
		userConversations = append(userConversations, userConversation)
	}

	return userConversations, nil

}

func (store *MessageStore) CreateMessage(senderId, receiverId, message string) error {
	_, err := store.db.Exec(
		`INSERT INTO messages (id, sender_id, receiver_id, message, created_at)
		 VALUES (?, ?, ?, ?, datetime('now'))`,
		uuid.New().String(), senderId, receiverId, message)
	if err != nil {
		return err
	}

	return nil
}
