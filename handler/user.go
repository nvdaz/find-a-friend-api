package handler

import (
	"net/http"

	"github.com/nvdaz/find-a-friend-api/db"

	"github.com/labstack/echo/v4"
)

func (handler *Handler) GetUser(c echo.Context) error {
	user, err := handler.userStore.GetUser(c.Param("id"))
	if err != nil {
		// TODO: 404 if user not found
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, user)
}

func (handler *Handler) GetAllUsers(c echo.Context) error {
	users, err := handler.userStore.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, users)
}

func (handler *Handler) CreateUser(c echo.Context) error {
	user := db.User{}
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
	}

	return c.NoContent(http.StatusCreated)
}

func (handler *Handler) UpdateUser(c echo.Context) error {
	user := db.PartialUser{}
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
	}

	return c.NoContent(http.StatusOK)
}
