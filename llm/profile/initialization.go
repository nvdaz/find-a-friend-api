package profile

import (
	"context"
	"fmt"
	"time"

	"github.com/nvdaz/find-a-friend-api/llm"
	"github.com/nvdaz/find-a-friend-api/model"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const (
	UserInterestsCount   = 10
	UserSkillsCount      = 5
	UserGoalsCount       = 3
	UserValuesCount      = 3
	UserExperiencesCount = 3
	UserHabitsCount      = 3
	UserHobbiesCount     = 5
)

func initializeInterests(questions string) ([]model.Interest, error) {
	system := fmt.Sprintf("Create a list of interests based on the provided chatbot questions. The list should contain %d specific interests. Provide a JSON object without any formatting containing the key 'interests', with the value being the list of interests. The interests should be objects with a key 'interest' containing the interest, a key 'level' containing the interest level on a scale of 0 to 1, and a key 'emoji' with a single, relevant emoji.", UserInterestsCount)

	result := struct {
		Interests []model.Interest `json:"interests"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelGpt4, questions, system, nil); err != nil {
		return nil, err
	}

	if len(result.Interests) == 0 {
		return nil, fmt.Errorf("no interests found")
	}

	return result.Interests, nil
}

func initializePersonality(questions string) (model.Personality, error) {
	system := "Let us play a guessing game. You are provided with a list of questions the user asked a chat bot. Guess the user's personality based on the Big Five (OCEAN) model, assigning scores from 0 to 5 for each trait, where 0 means the trait is not present and 5 signifies a strong presence. List the scores for Openness, Conscientiousness, Extroersion, Agreeableness, and Neuroticism. Start with analysis and use deductive reasoning to answer as precisely as possible. Then provide a JSON object in a JSON code block containing the keys 'openness', 'conscientiousness', 'extroversion', 'agreeableness', and 'neuroticism'."

	result := model.Personality{}
	if err := llm.GetResponseJson(&result, llm.ModelClaudeSonnet, questions, system, nil); err != nil {
		return model.Personality{}, err
	}

	return result, nil
}

func initializeSkills(questions string) ([]model.Skill, error) {
	system := fmt.Sprintf("Create a list of the user's skills based on the provided chatbot questions. The list should contain %d specific skills. Provide a JSON object without any formatting containing the key 'skills', with the value being the list of skills. The skills should be objects with a key 'skill' containing the skill and a key 'level' containing the skill level on a scale of 0 to 1.", UserSkillsCount)

	result := struct {
		Skills []model.Skill `json:"skills"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelClaudeSonnet, questions, system, nil); err != nil {
		return nil, err
	}

	if len(result.Skills) == 0 {
		return nil, fmt.Errorf("no skills found")
	}

	return result.Skills, nil
}

func initializeGoals(questions string) ([]model.Goal, error) {
	system := fmt.Sprintf("Let us play a guessing game. You are provided a list of chatbot questions a user asked. Guess the ambitions and goals that the user has. The list should contain %d specific goals. Describe goals in terse terms. Start with analysis and use deductive reasoning to answer as precisely as possible. Provide a JSON object containing the key 'goals', with the value being the list of goals. Each goal should be an object with a key 'goal' containing the goal and a key 'importance' containing the importance to the user on a scale from 0 to 1.", UserGoalsCount)

	result := struct {
		Goals []model.Goal `json:"goals"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelClaudeSonnet, questions, system, nil); err != nil {
		return nil, err
	}

	if len(result.Goals) == 0 {
		return nil, fmt.Errorf("no goals found")
	}

	return result.Goals, nil
}

func initializeValues(questions string) ([]model.CoreValue, error) {
	system := fmt.Sprintf("Create a list of the user's values and worldviews based on the provided chatbot questions. The list should contain %d specific values. Provide a JSON object in a JSON code block containing the key 'core_values', with the value being the list of values. Each value should be an object with a key 'value' containing the specific value and a key 'importance' containing the importance to the user on a scale from 0 to 1. Start with an in-depth analysis of the user's queries in an 'analysis' key.", UserValuesCount)

	result := struct {
		Values []model.CoreValue `json:"core_values"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelClaudeSonnet, questions, system, nil); err != nil {
		return nil, err
	}

	if len(result.Values) == 0 {
		return nil, fmt.Errorf("no core values found")
	}

	return result.Values, nil

}

func initializeDemographics(questions string) (model.Demographics, error) {
	system := "Let us play a guessing game. You are provided with a list of questions a user asked to a chatbot. Your task is to guess the user's demographic profile. Start with analysis and use deductive reasoning to answer as precisely as possible. Then, provide a valid JSON object in a JSON code block, with keys 'age', 'gender', 'location', 'occupation', 'highest_education', 'living_status', 'political_affiliation', 'religious_affiliation', 'nationality', 'spoken_languages' (list), and 'social_class'. Never reply with uncertainty; this is a game of deduction and analysis. Always provide a complete profile."
	result := model.Demographics{}
	if err := llm.GetResponseJson(&result, llm.ModelClaudeSonnet, questions, system, nil); err != nil {
		return model.Demographics{}, err
	}

	return result, nil
}

func initializeLivedExperiences(questions string) ([]string, error) {
	system := fmt.Sprintf("You are provided a list of questions a user asked to a chatbot. Create a list of %d specific lived experiences the user has had. Provide a JSON object without any formatting containing the key 'lived_experiences', with the value being the list of experiences.", UserExperiencesCount)

	result := struct {
		LivedExperiences []string `json:"lived_experiences"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelClaudeSonnet, questions, system, nil); err != nil {
		return nil, err
	}

	return result.LivedExperiences, nil
}

func initializeHabits(questions string) ([]string, error) {
	system := fmt.Sprintf("Create a list of %d specific habits based on the provided chatbot questions. Provide a JSON object without any formatting containing the key 'habits', with the value being the list of habits.", UserHabitsCount)

	result := struct {
		Habits []string `json:"habits"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelClaudeSonnet, questions, system, nil); err != nil {
		return nil, err
	}

	return result.Habits, nil
}

func initializeHobbies(questions string) ([]string, error) {
	system := fmt.Sprintf("Create a list of %d specific hobbies that the user does for fun in the form of verb phrases based on the provided chatbot questions. Provide a JSON object without any formatting containing the key 'hobbies', with the value being the list of hobbies.", UserHobbiesCount)

	result := struct {
		Hobbies []string `json:"hobbies"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelClaudeSonnet, questions, system, nil); err != nil {
		return nil, err
	}

	return result.Hobbies, nil

}

func initializeInterpersonalSkills(questions string) (model.InterpersonalSkills, error) {
	system := "Let us play a guessing game. You are provided with a list of questions a user asked to a chat bot. Guess their interpersonal skills. Start with analysis and use deductive reasoning to answer as precisely as possible. Then, provide a JSON object without any formatting containing the keys 'active_listening', 'teamwork', 'responsibility', 'dependability', 'leadership', 'motivation', 'flexibility', 'patience', and 'empathy'. Each key should have a value between 0 and 1, representing the strength of the skill."

	result := model.InterpersonalSkills{}
	if err := llm.GetResponseJson(&result, llm.ModelClaudeSonnet, questions, system, nil); err != nil {
		return model.InterpersonalSkills{}, err
	}

	return result, nil
}

func initializeExceptionalCircumstances(questions string) ([]string, error) {
	system := "Analyze the questions the user asked to identify any potential challenges or conditions they may have mentioned, such as disabilities or autism. Provide a JSON object, formatted properly, with the key 'exceptional_circumstances'. The value should be a list of these challenges, if any are mentioned. If no specific challenges are mentioned, the list should be empty. Start with an in-depth analysis of the user's queries in an 'analysis' key. "

	result := struct {
		ExceptionalCircumstances []string `json:"exceptional_circumstances"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelClaudeSonnet, questions, system, nil); err != nil {
		return nil, err
	}

	return result.ExceptionalCircumstances, nil

}

func initializeProfile(id string, questions string, conversations string) (*model.IntermediateProfile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	group, _ := errgroup.WithContext(ctx)

	sem := semaphore.NewWeighted(4)
	defer cancel()

	var interests []model.Interest
	var personality model.Personality
	var skills []model.Skill
	var goals []model.Goal
	var values []model.CoreValue
	var demographics model.Demographics
	var livedExperiences []string
	var interpersonalSkills model.InterpersonalSkills
	var habits []string
	var hobbies []string
	var exceptionalCircumstances []string
	var topics []model.Topic
	var conversationPersonality model.Personality
	var conversationInterpersonalSkills model.InterpersonalSkills
	var err error

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		interests, err = initializeInterests(questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		personality, err = initializePersonality(questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		skills, err = initializeSkills(questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		goals, err = initializeGoals(questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		values, err = initializeValues(questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		demographics, err = initializeDemographics(questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		livedExperiences, err = initializeLivedExperiences(questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		habits, err = initializeHabits(questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		interpersonalSkills, err = initializeInterpersonalSkills(questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		hobbies, err = initializeHobbies(questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		exceptionalCircumstances, err = initializeExceptionalCircumstances(questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)

		topics, err = GenerateTopicsFromConversations(conversations)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)

		conversationPersonality, err = GeneratePersonalityFromConversations(id, conversations)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)

		conversationInterpersonalSkills, err = GenerateInterpersonalSkillsFromConversations(id, conversations)
		return err
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	averagePersonality := model.Personality{
		Openness:          (personality.Openness + conversationPersonality.Openness) / 2,
		Conscientiousness: (personality.Conscientiousness + conversationPersonality.Conscientiousness) / 2,
		Extroversion:      (personality.Extroversion + conversationPersonality.Extroversion) / 2,
		Agreeableness:     (personality.Agreeableness + conversationPersonality.Agreeableness) / 2,
		Neuroticism:       (personality.Neuroticism + conversationPersonality.Neuroticism) / 2,
	}

	averageInterpersonalSkills := model.InterpersonalSkills{
		ActiveListening: (interpersonalSkills.ActiveListening + conversationInterpersonalSkills.ActiveListening) / 2,
		Teamwork:        (interpersonalSkills.Teamwork + conversationInterpersonalSkills.Teamwork) / 2,
		Responsibility:  (interpersonalSkills.Responsibility + conversationInterpersonalSkills.Responsibility) / 2,
		Dependability:   (interpersonalSkills.Dependability + conversationInterpersonalSkills.Dependability) / 2,
		Leadership:      (interpersonalSkills.Leadership + conversationInterpersonalSkills.Leadership) / 2,
		Motivation:      (interpersonalSkills.Motivation + conversationInterpersonalSkills.Motivation) / 2,
		Flexibility:     (interpersonalSkills.Flexibility + conversationInterpersonalSkills.Flexibility) / 2,
		Patience:        (interpersonalSkills.Patience + conversationInterpersonalSkills.Patience) / 2,
		Empathy:         (interpersonalSkills.Empathy + conversationInterpersonalSkills.Empathy) / 2,
	}

	return &model.IntermediateProfile{
		Interests:                interests,
		Personality:              averagePersonality,
		Skills:                   skills,
		Goals:                    goals,
		Values:                   values,
		Demographics:             demographics,
		LivedExperiences:         livedExperiences,
		Habits:                   habits,
		Hobbies:                  hobbies,
		InterpersonalSkills:      averageInterpersonalSkills,
		Topics:                   topics,
		ExceptionalCircumstances: exceptionalCircumstances,
	}, nil

}
