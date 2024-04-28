package service

import (
	"encoding/json"
	"fmt"

	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/llm/match"
	"github.com/nvdaz/find-a-friend-api/model"
)

type MatchService struct {
	UserService UserService
	MatchStore  db.MatchStore
}

func NewMatchService(userService UserService, matchStore db.MatchStore) MatchService {
	return MatchService{userService, matchStore}
}

func (service *MatchService) GetUserMatches(id string) ([]model.Match, error) {
	matches, err := service.MatchStore.GetUserMatches(id)
	if err != nil {
		return nil, err
	}

	convertedMatches := make([]model.Match, len(matches))
	for i, match := range matches {
		convertedMatches[i] = model.Match{
			UserId:  match.UserId,
			MatchId: match.MatchId,
			Reason:  match.Reason,
		}
	}

	return convertedMatches, nil
}

func (service *MatchService) GetAllNonMatchedUsers(id string) ([]model.User, error) {
	users, err := service.MatchStore.GetAllNonMatchedUsers(id)
	if err != nil {
		return nil, err
	}
	fmt.Println("Users", users)

	convertedUsers := make([]model.User, len(users))
	for i, user := range users {
		if user.Profile == nil {
			continue
		}

		profile := model.InternalProfile{}
		err = json.Unmarshal([]byte(*user.Profile), &profile)
		if err != nil {
			return nil, err
		}

		convertedUsers[i] = model.User{
			Id:      user.Id,
			Name:    user.Name,
			Profile: &profile,
		}

	}

	return convertedUsers, nil
}

func (service *MatchService) GenerateUserMatch(id string) (model.Match, error) {
	user, err := service.UserService.GetUser(id)
	if err != nil {
		return model.Match{}, err
	}

	otherUsers, err := service.GetAllNonMatchedUsers(id)
	if err != nil {
		return model.Match{}, err
	}

	fmt.Println("Other users", otherUsers)

	matchedUserId, err := match.GenerateMatch(*user, otherUsers)
	if err != nil {
		return model.Match{}, err
	}

	matchedUser := model.User{}
	for _, u := range otherUsers {
		if u.Id == *matchedUserId {
			matchedUser = u
			break
		}
	}

	if matchedUser.Id == "" {
		return model.Match{}, nil
	}

	firstMatchReason, err := match.ExplainMatchToUser(*user, matchedUser)
	if err != nil {
		return model.Match{}, err
	}

	secondMatchReason, err := match.ExplainMatchToUser(matchedUser, *user)
	if err != nil {
		return model.Match{}, err
	}

	err = service.MatchStore.CreateMatch(db.CreateMatch{
		UserId:  user.Id,
		MatchId: *matchedUserId,
		Reason:  firstMatchReason,
	}, db.CreateMatch{
		UserId:  *matchedUserId,
		MatchId: user.Id,
		Reason:  secondMatchReason,
	})
	if err != nil {
		return model.Match{}, err
	}

	return model.Match{
		UserId:  user.Id,
		MatchId: *matchedUserId,
		Reason:  firstMatchReason,
	}, nil
}
