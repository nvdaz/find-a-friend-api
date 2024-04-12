package db

import (
	"database/sql"
	"errors"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) UserStore {
	return UserStore{db}
}

type Interest struct {
	Interest  string  `json:"interest"`
	Intensity float64 `json:"intensity"`
	Skill     float64 `json:"skill"`
}

type Personality struct {
	Extraversion      float64 `json:"extraversion"`
	Agreeableness     float64 `json:"agreeableness"`
	Conscientiousness float64 `json:"conscientiousness"`
	Neuroticism       float64 `json:"neuroticism"`
	Openness          float64 `json:"openness"`
}

type User struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Personality Personality `json:"personality"`
	Interests   []Interest  `json:"interests"`
}

type PartialUser struct {
	Name        *string      `json:"name"`
	Personality *Personality `json:"personality"`
	Interests   *[]Interest  `json:"interests"`
}

var ErrNotFound = errors.New("not found")

func (store *UserStore) GetUser(id string) (*User, error) {
	user := User{}

	row := store.db.QueryRow(
		`SELECT u.id, u.name, p.extraversion, p.agreeableness, p.conscientiousness, p.neuroticism, p.openness
        	FROM users u
        	LEFT JOIN personalities p ON u.id = p.id
        	WHERE u.id = ?`,
		id)

	if err := row.Scan(&user.Id, &user.Name, &user.Personality.Extraversion,
		&user.Personality.Agreeableness, &user.Personality.Conscientiousness,
		&user.Personality.Neuroticism, &user.Personality.Openness); err != nil {

		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	rows, err := store.db.Query("SELECT interest, intensity, skill FROM interests WHERE user = ?", user.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interests []Interest
	for rows.Next() {
		var interest Interest
		if err := rows.Scan(&interest.Interest, &interest.Intensity, &interest.Skill); err != nil {
			return nil, err
		}
		interests = append(interests, interest)
	}

	user.Interests = interests

	return &user, nil
}

func (store *UserStore) GetAllUsers() ([]User, error) {
	rows, err := store.db.Query(
		`SELECT u.id, u.name, p.extraversion, p.agreeableness, p.conscientiousness, p.neuroticism, p.openness
			FROM users u
			LEFT JOIN personalities p ON u.id = p.id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.Id, &user.Name, &user.Personality.Extraversion,
			&user.Personality.Agreeableness, &user.Personality.Conscientiousness,
			&user.Personality.Neuroticism, &user.Personality.Openness); err != nil {
			return nil, err
		}

		rows, err := store.db.Query("SELECT interest, intensity, skill FROM interests WHERE user = ?", user.Id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var interests []Interest
		for rows.Next() {
			var interest Interest
			if err := rows.Scan(&interest.Interest, &interest.Intensity, &interest.Skill); err != nil {
				return nil, err
			}
			interests = append(interests, interest)
		}

		user.Interests = interests
		users = append(users, user)
	}

	return users, nil
}

func insertInterests(tx *sql.Tx, userId string, interests []Interest) error {
	if len(interests) == 0 {
		return nil
	}

	stmt, err := tx.Prepare("INSERT INTO interests (user, interest, intensity, skill) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, interest := range interests {
		if _, err := stmt.Exec(userId, interest.Interest, interest.Intensity, interest.Skill); err != nil {
			return err
		}
	}

	return nil
}

func (store *UserStore) CreateUser(user User) error {
	tx, err := store.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO users (id, name) VALUES (?, ?)", user.Id, user.Name)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO personalities (id, extraversion, agreeableness, conscientiousness, neuroticism, openness)
			VALUES (?, ?, ?, ?, ?, ?)`,
		user.Id, user.Personality.Extraversion, user.Personality.Agreeableness,
		user.Personality.Conscientiousness, user.Personality.Neuroticism, user.Personality.Openness,
	)

	if err != nil {
		return err
	}

	if err := insertInterests(tx, user.Id, user.Interests); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (store *UserStore) UpdateUser(id string, user PartialUser) error {
	tx, err := store.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT id FROM users WHERE id = ?", id)

	var test string
	if err := row.Scan(&test); err != nil {
		return err
	}

	if user.Name != nil {
		_, err = tx.Exec("UPDATE users SET name = ? WHERE id = ?", user.Name, id)
		if err != nil {
			return err
		}
	}

	if user.Personality != nil {
		_, err = tx.Exec(
			`UPDATE personalities
			SET extraversion = ?, agreeableness = ?, conscientiousness = ?, neuroticism = ?, openness = ?
			WHERE id = ?`,
			user.Personality.Extraversion, user.Personality.Agreeableness,
			user.Personality.Conscientiousness, user.Personality.Neuroticism, user.Personality.Openness, id,
		)
		if err != nil {
			return err
		}
	}

	if user.Interests != nil {
		_, err = tx.Exec("DELETE FROM interests WHERE user = ?", id)
		if err != nil {
			return err
		}

		if err := insertInterests(tx, id, *user.Interests); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
