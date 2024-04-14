package handler

import (
	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/service"
)

type Handler struct {
	userService              service.UserService
	userStore                *db.UserStore
	serviceConversationStore *db.ServiceConversationStore
}

func NewHandler(userService service.UserService, userStore *db.UserStore, serviceConversationStore *db.ServiceConversationStore) *Handler {
	return &Handler{userService, userStore, serviceConversationStore}
}
