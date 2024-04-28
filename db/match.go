package db

import (
	"database/sql"

	"github.com/google/uuid"
)

type MatchStore struct {
	db *sql.DB
}

func NewMatchStore(db *sql.DB) MatchStore {
	return MatchStore{db}
}

type Match struct {
	Id        string
	UserId    string
	MatchId   string
	Reason    string
	CreatedAt string
}

func (store *MatchStore) GetUserMatches(id string) ([]Match, error) {
	rows, err := store.db.Query(
		`SELECT user_id, match_id, reason, created_at
		 FROM matches
		 WHERE user_id = ?`,
		id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := []Match{}
	for rows.Next() {
		match := Match{}
		if err := rows.Scan(&match.UserId, &match.MatchId, &match.Reason, &match.CreatedAt); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}

type CreateMatch struct {
	UserId  string
	MatchId string
	Reason  string
}

func (store *MatchStore) CreateMatch(a, b CreateMatch) error {
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`INSERT INTO matches (id, user_id, match_id, reason, created_at)
		 VALUES (?, ?, ?, ?, datetime('now'))`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(uuid.New(), a.UserId, a.MatchId, a.Reason)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(uuid.New(), b.UserId, b.MatchId, b.Reason)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (store *MatchStore) GetAllNonMatchedUsers(id string) ([]User, error) {
	rows, err := store.db.Query(
		`SELECT id, name, updated_at, profile, generated_at
		 FROM users
		 WHERE id != ?
		 AND profile IS NOT NULL
		 AND id NOT IN (
			 SELECT match_id
			 FROM matches
			 WHERE user_id = ?
		 )`,
		id, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.Id, &user.Name, &user.UpdatedAt, &user.Profile, &user.GeneratedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
