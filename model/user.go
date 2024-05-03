package model

type Personality struct {
	Extroversion      float64 `json:"extroversion"`
	Agreeableness     float64 `json:"agreeableness"`
	Conscientiousness float64 `json:"conscientiousness"`
	Neuroticism       float64 `json:"neuroticism"`
	Openness          float64 `json:"openness"`
}

type User struct {
	Id      string           `json:"id"`
	Name    string           `json:"name"`
	Avatar  *string          `json:"avatar"`
	Profile *InternalProfile `json:"profile"`
}

type Interest struct {
	Interest string  `json:"interest"`
	Level    float64 `json:"level"`
	Emoji    string  `json:"emoji"`
}

type Skill struct {
	Skill string  `json:"skill"`
	Level float64 `json:"level"`
}

type Goal struct {
	Goal       string  `json:"goal"`
	Importance float64 `json:"importance"`
}

type CoreValue struct {
	Value      string  `json:"value"`
	Importance float64 `json:"importance"`
}

type MediaInterest struct {
	MediaInterest string  `json:"media_interest"`
	Level         float64 `json:"level"`
}

type Background struct {
	Occupation string `json:"occupation"`
	Education  string `json:"education"`
}

type Tag struct {
	Tag   string `json:"tag"`
	Emoji string `json:"emoji"`
}

type Topic struct {
	Topic string  `json:"topic"`
	Level float64 `json:"level"`
	Emoji string  `json:"emoji"`
}

type IntermediateProfile struct {
	Interests                []Interest          `json:"interests"`
	Personality              Personality         `json:"personality"`
	Skills                   []Skill             `json:"skills"`
	Goals                    []Goal              `json:"goals"`
	Values                   []CoreValue         `json:"values"`
	Demographics             Demographics        `json:"demographics"`
	LivedExperiences         []string            `json:"lived_experiences"`
	Habits                   []string            `json:"habits"`
	Hobbies                  []string            `json:"hobbies"`
	InterpersonalSkills      InterpersonalSkills `json:"interpersonal_skills"`
	ExceptionalCircumstances []string            `json:"exceptional_circumstances"`
	Topics                   []Topic             `json:"topics"`
}

type Demographics struct {
	AgeRange             string   `json:"age_range"`
	Gender               string   `json:"gender"`
	Occupation           string   `json:"occupation"`
	HighestEducation     string   `json:"highest_education"`
	LivingStatus         string   `json:"living_status"`
	PoliticalAffiliation string   `json:"political_affiliation"`
	ReligiousAffiliation string   `json:"religious_affiliation"`
	Nationality          string   `json:"nationality"`
	SpokenLanguages      []string `json:"spoken_languages"`
	SocialClass          string   `json:"social_class"`
}

type InterpersonalSkills struct {
	ActiveListening float64 `json:"active_listening"`
	Teamwork        float64 `json:"teamwork"`
	Responsibility  float64 `json:"responsibility"`
	Dependability   float64 `json:"dependability"`
	Leadership      float64 `json:"leadership"`
	Motivation      float64 `json:"motivation"`
	Flexibility     float64 `json:"flexibility"`
	Patience        float64 `json:"patience"`
	Empathy         float64 `json:"empathy"`
}

type ProfileFeatures struct {
	Summary      string   `json:"summary"`
	Tags         []Tag    `json:"tags"`
	Bio          string   `json:"bio"`
	KeyQuestions []string `json:"key_questions"`
	Subtitle     string   `json:"subtitle"`
	LookingFor   string   `json:"looking_for"`
}

type InternalProfile struct {
	Interests                []Interest          `json:"interests"`
	Personality              Personality         `json:"personality"`
	Skills                   []Skill             `json:"skills"`
	Goals                    []Goal              `json:"goals"`
	Values                   []CoreValue         `json:"values"`
	Demographics             Demographics        `json:"demographics"`
	LivedExperiences         []string            `json:"lived_experiences"`
	Habits                   []string            `json:"habits"`
	Hobbies                  []string            `json:"hobbies"`
	Topics                   []Topic             `json:"topics"`
	InterpersonalSkills      InterpersonalSkills `json:"interpersonal_skills"`
	ExceptionalCircumstances []string            `json:"exceptional_circumstances"`
	Summary                  string              `json:"summary"`
	Tags                     []Tag               `json:"tags"`
	Bio                      string              `json:"bio"`
	KeyQuestions             []string            `json:"key_questions"`
	Subtitle                 string              `json:"subtitle"`
	LookingFor               string              `json:"looking_for"`
}
