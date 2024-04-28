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
	if err := llm.GetResponseJson(&result, llm.ModelGpt3p5, questions, system, nil); err != nil {
		return nil, err
	}

	if len(result.Interests) == 0 {
		return nil, fmt.Errorf("no interests found")
	}

	return result.Interests, nil
}

func initializePersonality(questions string) (model.Personality, error) {
	system := "Assess the user's preliminary personality based on the Big Five (OCEAN) model, assigning scores from 0 to 5 for each trait, where 0 means the trait is not present and 5 signifies a strong presence. If it is not possible to determine a trait, provide an average value. List the scores for Openness, Conscientiousness, Extraversion, Agreeableness, and Neuroticism. These scores are only preliminary and need not be perfectly accurate. Provide a JSON object without any formatting containing the keys 'openness', 'conscientiousness', 'extraversion', 'agreeableness', and 'neuroticism'."

	result := model.Personality{}
	if err := llm.GetResponseJson(&result, llm.ModelGpt3p5, questions, system, nil); err != nil {
		return model.Personality{}, err
	}

	return result, nil
}

func initializeSkills(questions string) ([]model.Skill, error) {
	system := fmt.Sprintf("Create a list of the user's skills based on the provided chatbot questions. The list should contain %d specific skills. Provide a JSON object without any formatting containing the key 'skills', with the value being the list of skills. The skills should be objects with a key 'skill' containing the skill and a key 'level' containing the skill level on a scale of 0 to 1.", UserSkillsCount)

	result := struct {
		Skills []model.Skill `json:"skills"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelGpt3p5, questions, system, nil); err != nil {
		return nil, err
	}

	if len(result.Skills) == 0 {
		return nil, fmt.Errorf("no skills found")
	}

	return result.Skills, nil
}

func initializeGoals(questions string) ([]model.Goal, error) {
	system := fmt.Sprintf("Create a list of ambitions and goals that the user has based on the provided chatbot questions. The list should contain %d specific goals. Describe goals in terse terms. Provide a JSON object without any formatting containing the key 'goals', with the value being the list of goals. Each goal should be an object with a key 'goal' containing the goal and a key 'importance' containing the importance to the user on a scale from 0 to 1.", UserGoalsCount)

	result := struct {
		Goals []model.Goal `json:"goals"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelGpt3p5, questions, system, nil); err != nil {
		return nil, err
	}

	if len(result.Goals) == 0 {
		return nil, fmt.Errorf("no goals found")
	}

	return result.Goals, nil
}

func initializeValues(questions string) ([]model.CoreValue, error) {
	system := fmt.Sprintf("Create a list of the user's values and worldviews based on the provided chatbot questions. The list should contain %d specific values. Provide a JSON object without any formatting containing the key 'core_values', with the value being the list of values. Each value should be an object with a key 'value' containing the specific value and a key 'importance' containing the importance to the user on a scale from 0 to 1.", UserValuesCount)

	result := struct {
		Values []model.CoreValue `json:"core_values"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelGpt3p5, questions, system, nil); err != nil {
		return nil, err
	}

	if len(result.Values) == 0 {
		return nil, fmt.Errorf("no core values found")
	}

	return result.Values, nil

}

func initializeDemographics(questions string) (model.Demographics, error) {
	system := "Analyze user-asked questions to deduce their demographic profile. Provide a valid JSON object without formatting, with keys 'age_range', 'gender', 'occupation', 'highest_education', 'living_status', 'political_affiliation', 'religious_affiliation', 'nationality', 'spoken_languages' (list), and 'social_class'. Start with an in-depth analysis of the user's queries in an 'analysis' key. Never reply with uncertainty; always provide a best answer."
	result := model.Demographics{}
	if err := llm.GetResponseJson(&result, llm.ModelGpt4, questions, system, nil); err != nil {
		return model.Demographics{}, err
	}

	return result, nil
}

func initializeLivedExperiences(questions string) ([]string, error) {
	system := fmt.Sprintf("Create a list of %d specific lived experiences based on the provided chatbot questions. Provide a JSON object without any formatting containing the key 'lived_experiences', with the value being the list of experiences.", UserExperiencesCount)

	result := struct {
		LivedExperiences []string `json:"lived_experiences"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelGpt3p5, questions, system, nil); err != nil {
		return nil, err
	}

	return result.LivedExperiences, nil
}

func initializeHabits(questions string) ([]string, error) {
	system := fmt.Sprintf("Create a list of %d specific habits based on the provided chatbot questions. Provide a JSON object without any formatting containing the key 'habits', with the value being the list of habits.", UserHabitsCount)

	result := struct {
		Habits []string `json:"habits"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelGpt3p5, questions, system, nil); err != nil {
		return nil, err
	}

	return result.Habits, nil
}

func initializeHobbies(questions string) ([]string, error) {
	system := fmt.Sprintf("Create a list of %d specific hobbies that the user does for fun in the form of verb phrases based on the provided chatbot questions. Provide a JSON object without any formatting containing the key 'hobbies', with the value being the list of hobbies.", UserHobbiesCount)

	result := struct {
		Hobbies []string `json:"hobbies"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelGpt3p5, questions, system, nil); err != nil {
		return nil, err
	}

	return result.Hobbies, nil

}

func initializeInterpersonalSkills(questions string) (model.InterpersonalSkills, error) {
	system := "Analyze the user's responses to determine their interpersonal skills. Provide a JSON object without any formatting containing the keys 'active_listening', 'teamwork', 'responsibility', 'dependability', 'leadership', 'motivation', 'flexibility', 'patience', and 'empathy'. Each key should have a value between 0 and 1, representing the strength of the skill."

	result := model.InterpersonalSkills{}
	if err := llm.GetResponseJson(&result, llm.ModelGpt4, questions, system, nil); err != nil {
		return model.InterpersonalSkills{}, err
	}

	return result, nil
}

func initializeExceptionalCircumstances(questions string) ([]string, error) {
	system := "Analyze the questions the user asked to identify any potential challenges or conditions they may have mentioned, such as disabilities or autism. Provide a JSON object, formatted properly, with the key 'exceptional_circumstances'. The value should be a list of these challenges, if any are mentioned. If no specific challenges are mentioned, the list should be empty."

	result := struct {
		ExceptionalCircumstances []string `json:"exceptional_circumstances"`
	}{}
	if err := llm.GetResponseJson(&result, llm.ModelGpt3p5, questions, system, nil); err != nil {
		return nil, err
	}

	return result.ExceptionalCircumstances, nil

}

func initializeProfile(questions string) (*model.IntermediateProfile, error) {
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
	var err error

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)
		interests, err = initializeInterests(questions)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)
		personality, err = initializePersonality(questions)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)
		skills, err = initializeSkills(questions)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)
		goals, err = initializeGoals(questions)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)
		values, err = initializeValues(questions)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)
		demographics, err = initializeDemographics(questions)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)
		livedExperiences, err = initializeLivedExperiences(questions)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)
		habits, err = initializeHabits(questions)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)
		interpersonalSkills, err = initializeInterpersonalSkills(questions)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)
		hobbies, err = initializeHobbies(questions)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)
		exceptionalCircumstances, err = initializeExceptionalCircumstances(questions)
		return err
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return &model.IntermediateProfile{
		Interests:                interests,
		Personality:              personality,
		Skills:                   skills,
		Goals:                    goals,
		Values:                   values,
		Demographics:             demographics,
		LivedExperiences:         livedExperiences,
		Habits:                   habits,
		Hobbies:                  hobbies,
		InterpersonalSkills:      interpersonalSkills,
		ExceptionalCircumstances: exceptionalCircumstances,
	}, nil

}
