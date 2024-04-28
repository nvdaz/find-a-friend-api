package handler

import (
	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/service"
)

type Handler struct {
	userService              service.UserService
	matchService             service.MatchService
	serviceConversationStore *db.ServiceConversationStore
}

func NewHandler(userService service.UserService, matchService service.MatchService, serviceConversationStore *db.ServiceConversationStore) *Handler {
	return &Handler{userService, matchService, serviceConversationStore}
}
