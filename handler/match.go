package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (handler *Handler) GetUserMatches(c echo.Context) error {
	id := c.Param("id")
	matches, err := handler.matchService.GetMatchedUsers(id)
	if err != nil {
		fmt.Println("Error getting user matches", err)
		return echo.NewHTTPError(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, matches)
}

func (handler *Handler) GenerateUserMatch(c echo.Context) error {
	id := c.Param("id")
	match, err := handler.matchService.GenerateUserMatch(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, match)
}

func (handler *Handler) GetMatch(c echo.Context) error {
	id := c.Param("id")
	match, err := handler.matchService.GetMatch(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, match)
}
