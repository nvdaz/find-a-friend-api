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

type User struct {
	Id          string
	Name        string
	UpdatedAt   string
	Profile     *string
	GeneratedAt *string
}

var ErrUserNotFound = errors.New("user not found")

func (store *UserStore) GetUser(id string) (*User, error) {
	user := User{}

	row := store.db.QueryRow("SELECT id, name, updated_at, profile, generated_at FROM users WHERE id = ?", id)

	if err := row.Scan(&user.Id, &user.Name, &user.UpdatedAt, &user.Profile, &user.GeneratedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

type CreateUser struct {
	Id   string
	Name string
}

func (store *UserStore) CreateUser(createUser CreateUser) error {
	_, err := store.db.Exec("INSERT INTO users (id, name, updated_at) VALUES (?, ?, datetime('now'))", createUser.Id, createUser.Name)

	return err
}

func (store *UserStore) MarkUserAsUpdated(id string) error {
	_, err := store.db.Exec("UPDATE users SET updated_at = datetime('now') WHERE id = ?", id)

	return err
}

func (store *UserStore) UpdateUser(user User) error {
	_, err := store.db.Exec("UPDATE users SET name = ?, updated_at = ? WHERE id = ?", user.Name, user.UpdatedAt, user.Id)

	return err
}

func (store *UserStore) UpdateUserProfile(id string, profile string) error {
	_, err := store.db.Exec("UPDATE users SET profile = ?, generated_at = datetime('now') WHERE id = ?", profile, id)

	return err
}

func (store *UserStore) GetAllUsers() ([]User, error) {
	rows, err := store.db.Query("SELECT id, name, updated_at, profile, generated_at FROM users")
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
