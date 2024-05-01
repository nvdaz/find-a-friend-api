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
	OtherId   string
	Reason    string
	CreatedAt string
}

func (store *MatchStore) GetUserMatches(id string) ([]Match, error) {
	rows, err := store.db.Query(
		`SELECT id, user_id, other_id, reason, created_at
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
		if err := rows.Scan(&match.Id, &match.UserId, &match.OtherId, &match.Reason, &match.CreatedAt); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}

type CreateMatch struct {
	UserId  string
	OtherId string
	Reason  string
}

func (store *MatchStore) CreateMatch(a, b CreateMatch) (*string, error) {
	tx, err := store.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`INSERT INTO matches (id, user_id, other_id, reason, created_at)
		 VALUES (?, ?, ?, ?, datetime('now'))`)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()

	_, err = stmt.Exec(id, a.UserId, a.OtherId, a.Reason)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(uuid.New(), b.UserId, b.OtherId, b.Reason)
	if err != nil {
		return nil, err
	}

	return &id, tx.Commit()
}

func (store *MatchStore) GetAllNonMatchedUsers(id string) ([]User, error) {
	rows, err := store.db.Query(
		`SELECT id, name, updated_at, profile, generated_at
		 FROM users
		 WHERE id != ?
		 AND profile IS NOT NULL
		 AND id NOT IN (
			 SELECT other_id
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

func (store *MatchStore) GetMatch(id string) (Match, error) {
	row := store.db.QueryRow(
		`SELECT id, user_id, other_id, reason, created_at
		 FROM matches
		 WHERE id = ?`,
		id)

	match := Match{}
	if err := row.Scan(&match.Id, &match.UserId, &match.OtherId, &match.Reason, &match.CreatedAt); err != nil {
		return Match{}, err
	}

	return match, nil
}

func (store *MatchStore) GetMatchedUsers(id string) ([]User, error) {
	rows, err := store.db.Query(
		`SELECT id, name, avatar, updated_at, profile, generated_at
		 FROM users
		 WHERE id IN (
			 SELECT other_id
			 FROM matches
			 WHERE user_id = ?
		 )`,
		id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.Id, &user.Name, &user.Avatar, &user.UpdatedAt, &user.Profile, &user.GeneratedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
