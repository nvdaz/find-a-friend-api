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
	UserTagsCount         = 8
	UserKeyQuestionsCount = 3
)

func generateUserBio(user *model.NonSecretIntermediateProfile) (string, error) {
	profileString, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	result := struct {
		Bio string `json:"bio"`
	}{}

	err = llm.GetResponseJson(&result, llm.ModelGpt4, string(profileString), "Create a short, passionate introductory biography from the perspective of the provided user using personal pronouns. Include a brief description of their personality and interests. The biography should be a single paragraph, no more than 200 words in length. Provide a JSON object without any formatting containing two keys: 'bio', with the value being the biography.", nil)
	if err != nil {
		return "", err
	}

	return result.Bio, nil
}

func generateUserKeyQuestions(user *model.NonSecretIntermediateProfile) ([]string, error) {
	profileString, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	result := struct {
		Questions []string `json:"key_questions"`
	}{}

	err = llm.GetResponseJson(&result, llm.ModelGpt4, string(profileString), fmt.Sprintf("Create a list of %d key questions that the user has already asked the chat bot that are representative of their interests and selected to spark conversation. Provide a JSON object without any formatting containing a single key: 'key_questions', with the value being a list of the questions.", UserKeyQuestionsCount), nil)
	if err != nil {
		return nil, err
	}

	return result.Questions, nil
}

func generateUserTags(user *model.NonSecretIntermediateProfile) ([]model.Tag, error) {
	profileString, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	result := struct {
		Tags []model.Tag `json:"tags"`
	}{}

	err = llm.GetResponseJson(&result, llm.ModelGpt4, string(profileString), fmt.Sprintf("Create a list of %d short tags that describe the user. The tags should be representative of their interests and personality. Provide a JSON object without any formatting containing a single key: 'tags', with the value being a list of tags. Each tag should have a key 'tag' with the tag name and a key 'emoji' with an emoji to accompany it.", UserTagsCount), nil)
	if err != nil {
		return nil, err
	}

	return result.Tags, nil
}

func generateUserSummary(user *model.NonSecretIntermediateProfile) (string, error) {
	profileString, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	result := struct {
		Summary string `json:"summary"`
	}{}

	err = llm.GetResponseJson(&result, llm.ModelGpt4, string(profileString), "Create a short summary of the user's profile including only the most important information about them. The summary should be no more than 200 words in length. Provide a JSON object without any formatting containing a single key: 'summary', with the value being the summary.", nil)
	if err != nil {
		return "", err
	}

	return result.Summary, nil

}

func generateUserFeatures(user *model.NonSecretIntermediateProfile) (*model.ProfileFeatures, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	group, _ := errgroup.WithContext(ctx)

	sem := semaphore.NewWeighted(4)
	defer cancel()

	var summary string
	var tags []model.Tag
	var bio string
	var keyQuestions []string
	var err error

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)

		summary, err = generateUserSummary(user)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)

		tags, err = generateUserTags(user)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)

		bio, err = generateUserBio(user)
		return err
	})

	group.Go(func() error {
		sem.Acquire(ctx, 1)
		defer sem.Release(1)

		keyQuestions, err = generateUserKeyQuestions(user)
		return err
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return &model.ProfileFeatures{
		Summary: summary,
		Tags:         tags,
		Bio:          bio,
		KeyQuestions: keyQuestions,
	}, nil
}
