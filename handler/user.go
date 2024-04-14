package handler

import (
	"log"
	"net/http"

	"github.com/nvdaz/find-a-friend-api/db"

	"github.com/labstack/echo/v4"
)

func (handler *Handler) GetUser(c echo.Context) error {
	user, err := handler.userService.GetUser(c.Param("id"))
	if err != nil {
		if err == db.ErrUserNotFound {
			return c.JSON(http.StatusNotFound, nil)
		}
		log.Println("Error getting user:", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, user)
}

func (handler *Handler) CreateUser(c echo.Context) error {
	user := db.User{}
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
	}

	return c.NoContent(http.StatusCreated)
}

func (handler *Handler) UpdateUser(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}
