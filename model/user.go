package model

type Personality struct {
	Extroversion      float64 `json:"extroversion"`
	Agreeableness     float64 `json:"agreeableness"`
	Conscientiousness float64 `json:"conscientiousness"`
	Neuroticism       float64 `json:"neuroticism"`
	Openness          float64 `json:"openness"`
}

type Interest struct {
	Interest  string  `json:"interest"`
	Intensity float64 `json:"intensity"`
	Skill     float64 `json:"skill"`
}

type User struct {
	Id           string      `json:"id"`
	Name         string      `json:"name"`
	Bio          string      `json:"bio"`
	Personality  Personality `json:"personality"`
	Interests    []Interest  `json:"interests"`
	KeyQuestions []string    `json:"key_questions"`
}

type Match struct {
	Id       string  `json:"id"`
	UserId   string  `json:"user_id"`
	Affinity float64 `json:"affinity"`
	Reason   string  `json:"reason"`
}
