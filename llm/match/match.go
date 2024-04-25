package match

import (
	"context"
	"fmt"
	"time"

	"github.com/nvdaz/find-a-friend-api/model"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func GenerateMatch(user model.User, users []model.User) (*model.Match, error) {
	candidates, err := GenerateCandidateMatches(user, users)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	group, _ := errgroup.WithContext(ctx)

	sem := semaphore.NewWeighted(4)
	defer cancel()

	explanations := make(map[string]string)
	for _, candidateId := range candidates {
		var candidateUser *model.User
		for _, u := range users {
			if u.Id == candidateId {
				candidateUser = &u
				break
			}
		}

		if candidateUser == nil {
			continue
		}

		group.Go(func() error {
			sem.Acquire(ctx, 1)
			defer sem.Release(1)

			explanation, err := ExplainMatch(user, *candidateUser)
			if err != nil {
				return err
			}

			explanations[candidateId] = explanation
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	bestMatchId, err := DecideBestMatch(explanations)
	if err != nil {
		return nil, err
	}

	var bestUser *model.User

	for _, u := range users {
		if u.Id == bestMatchId {
			bestUser = &u
		}
	}

	if bestUser == nil {
		return nil, fmt.Errorf("got invalid id")
	}

	reason, err := ExplainMatchToUser(user, *bestUser)
	if err != nil {
		return nil, err
	}

	return &model.Match{
		User:   *bestUser,
		Reason: reason,
	}, nil
}
