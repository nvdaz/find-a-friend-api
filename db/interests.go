package db

import (
	"database/sql"
)

type InterestsStore struct {
	db *sql.DB
}

func NewInterestsStore(db *sql.DB) InterestsStore {
	return InterestsStore{db}
}

type Interest struct {
	Interest  string
	Intensity float64
	Skill     float64
}

func (store *InterestsStore) GetUserInterests(user_id string) ([]Interest, error) {
	rows, err := store.db.Query("SELECT interest, intensity, skill FROM interests WHERE user_id = ?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	interests := []Interest{}

	for rows.Next() {
		interest := Interest{}
		if err := rows.Scan(&interest.Interest, &interest.Intensity, &interest.Skill); err != nil {
			return nil, err
		}
		interests = append(interests, interest)
	}

	return interests, nil
}

func (store *InterestsStore) InsertUserInterests(user_id string, interests []Interest) error {
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM interests WHERE user_id = ?", user_id)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO interests (user_id, interest, intensity, skill) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, interest := range interests {
		if _, err := stmt.Exec(user_id, interest.Interest, interest.Intensity, interest.Skill); err != nil {
			return err
		}
	}

	return tx.Commit()
}
