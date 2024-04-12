package handler

import (
	"github.com/nvdaz/find-a-friend-api/db"
)

type Handler struct {
	userStore *db.UserStore
}

func NewHandler(userStore *db.UserStore) *Handler {
	return &Handler{userStore}
}
