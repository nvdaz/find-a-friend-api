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
	Extroversion      float64 `json:"extroversion"`
	Agreeableness     float64 `json:"agreeableness"`
	Conscientiousness float64 `json:"conscientiousness"`
	Neuroticism       float64 `json:"neuroticism"`
	Openness          float64 `json:"openness"`
}

type UserProfile struct {
	Id          string      `json:"id"`
	Bio         string      `json:"bio"`
	Personality Personality `json:"personality"`
	UpdatedAt   string      `json:"updated_at"`
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
