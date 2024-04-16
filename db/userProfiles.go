package db

import (
	"database/sql"
	"errors"
)

type UserProfilesStore struct {
	db *sql.DB
}

func NewUserProfilesStore(db *sql.DB) UserProfilesStore {
	return UserProfilesStore{db}
}

type Personality struct {
	Extroversion      float64
	Agreeableness     float64
	Conscientiousness float64
	Neuroticism       float64
	Openness          float64
}

type UserProfile struct {
	Id          string
	Bio         string
	Personality Personality
	UpdatedAt   string
}

var ErrUserProfileNotFound = errors.New("user profile not found")

func (store *UserProfilesStore) GetUserProfile(id string) (*UserProfile, error) {
	userProfile := UserProfile{}

	row := store.db.QueryRow("SELECT id, bio, extroversion, agreeableness, conscientiousness, neuroticism, openness, updated_at FROM user_profiles WHERE id = ?", id)

	if err := row.Scan(&userProfile.Id, &userProfile.Bio,
		&userProfile.Personality.Extroversion, &userProfile.Personality.Agreeableness,
		&userProfile.Personality.Conscientiousness, &userProfile.Personality.Neuroticism,
		&userProfile.Personality.Openness, &userProfile.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, ErrUserProfileNotFound
		}
		return nil, err
	}

	return &userProfile, nil
}

func (store *UserProfilesStore) InsertUserProfile(userProfile UserProfile) error {
	_, err := store.db.Exec(`
		INSERT INTO user_profiles (id, bio, extroversion, agreeableness, conscientiousness, neuroticism, openness, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'))
		ON CONFLICT (id) DO UPDATE SET
			bio = excluded.bio,
			extroversion = excluded.extroversion,
			agreeableness = excluded.agreeableness,
			conscientiousness = excluded.conscientiousness,
			neuroticism = excluded.neuroticism,
			openness = excluded.openness,
			updated_at = datetime('now')`,
		userProfile.Id, userProfile.Bio, userProfile.Personality.Extroversion,
		userProfile.Personality.Agreeableness, userProfile.Personality.Conscientiousness,
		userProfile.Personality.Neuroticism, userProfile.Personality.Openness)

	return err
}

func (store *UserProfilesStore) GetAllUserProfiles() ([]UserProfile, error) {
	rows, err := store.db.Query("SELECT id, bio, extroversion, agreeableness, conscientiousness, neuroticism, openness, updated_at FROM user_profiles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userProfiles := []UserProfile{}
	for rows.Next() {
		userProfile := UserProfile{}
		if err := rows.Scan(&userProfile.Id, &userProfile.Bio,
			&userProfile.Personality.Extroversion, &userProfile.Personality.Agreeableness,
			&userProfile.Personality.Conscientiousness, &userProfile.Personality.Neuroticism,
			&userProfile.Personality.Openness, &userProfile.UpdatedAt); err != nil {
			return nil, err
		}
		userProfiles = append(userProfiles, userProfile)
	}

	return userProfiles, nil
}
