package handler

import (
	"fmt"
	"net/http"

	"github.com/nvdaz/find-a-friend-api/db"

	"github.com/labstack/echo/v4"
)

func (handler *Handler) GetUser(c echo.Context) error {
	user, err := handler.userService.GetUser(c.Param("id"))
	if err != nil {
		fmt.Println("Error getting user", err)
		if err == db.ErrUserNotFound {
			return c.JSON(http.StatusNotFound, nil)
		}

		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, user)
}

func (handler *Handler) GetAllUsers(c echo.Context) error {
	users, err := handler.userService.GetAllUsers()
	if err != nil {
		fmt.Println("Error getting users", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, users)
}

func (handler *Handler) CreateUser(c echo.Context) error {
	user := db.CreateUser{}
	if err := c.Bind(&user); err != nil {
		fmt.Println("Error parsing", err)
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
	}

	if err := handler.userService.CreateUser(user); err != nil {
		fmt.Println("Error creating", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error creating user")
	}

	return c.NoContent(http.StatusCreated)
}

func (handler *Handler) UpdateUser(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (handler *Handler) GetUserMatches(c echo.Context) error {
	matches, err := handler.userService.GetBestMatch(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting matches")
	}

	return c.JSON(http.StatusOK, matches)
}
