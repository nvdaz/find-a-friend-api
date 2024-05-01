package handler

import (
	"github.com/nvdaz/find-a-friend-api/service"
)

type Handler struct {
	userService    service.UserService
	matchService   service.MatchService
	messageService service.MessageService
}

func NewHandler(userService service.UserService, matchService service.MatchService, messageService service.MessageService) *Handler {
	return &Handler{userService, matchService, messageService}
}
