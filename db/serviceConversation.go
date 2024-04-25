package db

import (
	"database/sql"
	"fmt"
)

type ServiceConversationStore struct {
	db *sql.DB
}

func NewServiceConversationStore(db *sql.DB) ServiceConversationStore {
	return ServiceConversationStore{db}
}

type ServiceConversation struct {
	Id        string
	UserId    string
	Question  string
	Answer    string
	CreatedAt string
	IsKey     bool
}

func (store *ServiceConversationStore) GetRecentServiceConversations(userId string, limit int) ([]ServiceConversation, error) {
	rows, err := store.db.Query(
		`SELECT id, user_id, question, answer, created_at, is_key
		 FROM service_conversations
		 WHERE user_id = ?
		 ORDER BY created_at DESC
		 LIMIT ?`,
		userId, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	serviceConversations := []ServiceConversation{}
	for rows.Next() {
		serviceConversation := ServiceConversation{}
		if err := rows.Scan(&serviceConversation.Id, &serviceConversation.UserId, &serviceConversation.Question,
			&serviceConversation.Answer, &serviceConversation.CreatedAt, &serviceConversation.IsKey); err != nil {
			return nil, err
		}
		serviceConversations = append(serviceConversations, serviceConversation)
	}

	return serviceConversations, nil
}

func (store *ServiceConversationStore) CreateServiceConversation(serviceConversation ServiceConversation) error {
	_, err := store.db.Exec(
		`INSERT INTO service_conversations (id, user_id, question, answer, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		serviceConversation.Id, serviceConversation.UserId, serviceConversation.Question, serviceConversation.Answer, serviceConversation.CreatedAt)

	return err
}

func (store *ServiceConversationStore) CreateServiceConversations(serviceConversations []ServiceConversation) error {
	tx, err := store.db.Begin()
	if err != nil {
		fmt.Println("Error creating transaction", err)
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`INSERT INTO service_conversations (id, user_id, question, answer, created_at)
		 VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, serviceConversation := range serviceConversations {
		if _, err := stmt.Exec(serviceConversation.Id, serviceConversation.UserId, serviceConversation.Question, serviceConversation.Answer, serviceConversation.CreatedAt); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (store *ServiceConversationStore) UpdateKeyServiceConversations(userId string, ids []string) error {
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec("UPDATE service_conversations SET is_key = false WHERE user_id = ?", userId)

	stmt, err := tx.Prepare("UPDATE service_conversations SET is_key = true WHERE user_id = ? AND id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, id := range ids {
		if _, err := stmt.Exec(userId, id); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (store *ServiceConversationStore) GetKeyServiceConversations(userId string) ([]ServiceConversation, error) {
	rows, err := store.db.Query("SELECT id, question, answer, created_at FROM service_conversations WHERE user_id = ? AND is_key = true", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	serviceConversations := []ServiceConversation{}
	for rows.Next() {
		serviceConversation := ServiceConversation{UserId: userId, IsKey: true}
		if err := rows.Scan(&serviceConversation.Id, &serviceConversation.Question, &serviceConversation.Answer, &serviceConversation.CreatedAt); err != nil {
			return nil, err
		}
		serviceConversations = append(serviceConversations, serviceConversation)
	}

	return serviceConversations, nil
}
