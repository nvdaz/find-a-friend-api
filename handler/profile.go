package handler

import (
	"net/http"

	"github.com/nvdaz/find-a-friend-api/llm"

	"github.com/labstack/echo/v4"
)

func (handler *Handler) GetUserProfile(c echo.Context) error {
	user, err := handler.userStore.GetUser(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	bio, err := llm.GenerateUserBio(*user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, bio)
}
