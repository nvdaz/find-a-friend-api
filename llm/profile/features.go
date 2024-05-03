package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nvdaz/find-a-friend-api/llm"
	"github.com/nvdaz/find-a-friend-api/model"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

const (
	UserTagsCount         = 4
	UserKeyQuestionsCount = 3
)

func generateUserBio(user model.IntermediateProfile) (string, error) {
	profileString, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	result := struct {
		Bio string `json:"bio"`
	}{}

	err = llm.GetResponseJson(&result, llm.ModelGpt4, string(profileString), "Create a short, passionate introductory biography in a casual, friendly tone from the perspective of the provided user using personal pronouns. Include a brief description of their personality and interests. The biography should be a single paragraph, no more than 120 words in length. Provide a JSON object without any formatting containing two keys: 'bio', with the value being the biography.", nil)
	if err != nil {
		return "", err
	}

	return result.Bio, nil
}

func generateUserKeyQuestions(user model.IntermediateProfile, questions string) ([]string, error) {
	prompt := struct {
		User      model.IntermediateProfile `json:"user"`
		Questions string                    `json:"questions"`
	}{
		User:      user,
		Questions: questions,
	}
	data, err := json.Marshal(prompt)
	if err != nil {
		return nil, err
	}

	result := struct {
		Questions []string `json:"key_questions"`
	}{}

	err = llm.GetResponseJson(&result, llm.ModelClaudeSonnet, string(data), fmt.Sprintf("Create a list of %d key questions that the user has already asked the chat bot that are representative of their interests and selected to spark conversation. Provide a JSON object without any formatting containing a single key: 'key_questions', with the value being a list of the questions.", UserKeyQuestionsCount), nil)
	if err != nil {
		return nil, err
	}

	return result.Questions, nil
}

func generateUserTags(user model.IntermediateProfile) ([]model.Tag, error) {
	profileString, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	result := struct {
		Tags []model.Tag `json:"tags"`
	}{}

	err = llm.GetResponseJson(&result, llm.ModelGpt4, string(profileString), fmt.Sprintf("Create a list of %d short tags that describe the user. The tags should be representative of who they are, but not restating what is already given (for example: analytical thinker, in college, ethical innovator). Provide a JSON object without any formatting containing a single key: 'tags', with the value being a list of tags. Each tag should have a key 'tag' with the tag name and a key 'emoji' with a single emoji to accompany it.", UserTagsCount), nil)
	if err != nil {
		return nil, err
	}

	return result.Tags, nil
}

func generateUserSummary(user model.IntermediateProfile) (string, error) {
	profileString, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	result := struct {
		Summary string `json:"summary"`
	}{}

	err = llm.GetResponseJson(&result, llm.ModelGpt4, string(profileString), "Create an in-depth summary of the user's profile including only the most important information about them. The summary should be no more than 120 words in length. Provide a JSON object without any formatting containing a single key: 'summary', with the value being the summary.", nil)
	if err != nil {
		return "", err
	}

	return result.Summary, nil
}

func generateUserSubtitle(user model.IntermediateProfile) (string, error) {
	profileString, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	result := struct {
		Subtitle string `json:"subtitle"`
	}{}

	err = llm.GetResponseJson(&result, llm.ModelGpt4, string(profileString), "Create a 2-6 word creative subtitle in a casual, friendly tone to go under the user's name under their profile that captures the essence of their personality. Be as unique and creative as possible. Dive into what cannot be immediately seen just by their profile. Provide a JSON object without any formatting containing a single key: 'subtitle', with the value being the subtitle", nil)
	if err != nil {
		return "", err
	}

	return result.Subtitle, nil
}

func generateUserLookingFor(user model.IntermediateProfile) (string, error) {
	profileString, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	result := struct {
		LookingFor string `json:"looking_for"`
	}{}

	err = llm.GetResponseJson(&result, llm.ModelGpt4, string(profileString), "Create a short, creative description (about 10-15 words) that expresses the kind of friend the user is looing for (for example: Like-minded girlfriends to share a love of books and coffee). Be as unique and creative as possible. Dive into what cannot be immediately seen just by their profile. Provide a JSON object without any formatting containing a single key: 'looking_for', with the value being the description.", nil)
	if err != nil {
		return "", err
	}

	return result.LookingFor, nil

}

func generateUserFeatures(user model.IntermediateProfile, questions string) (*model.ProfileFeatures, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	group, _ := errgroup.WithContext(ctx)

	sem := semaphore.NewWeighted(4)
	defer cancel()

	var summary string
	var tags []model.Tag
	var bio string
	var keyQuestions []string
	var subtitle string
	var lookingFor string
	var err error

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)

		summary, err = generateUserSummary(user)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)

		tags, err = generateUserTags(user)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)

		bio, err = generateUserBio(user)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)

		keyQuestions, err = generateUserKeyQuestions(user, questions)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)

		subtitle, err = generateUserSubtitle(user)
		return err
	})

	group.Go(func() error {
		if err = sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)

		lookingFor, err = generateUserLookingFor(user)
		return err
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return &model.ProfileFeatures{
		Summary:      summary,
		Tags:         tags,
		Bio:          bio,
		KeyQuestions: keyQuestions,
		Subtitle:     subtitle,
		LookingFor:   lookingFor,
	}, nil
}
