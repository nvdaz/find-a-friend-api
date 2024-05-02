package handler

import (
	"fmt"
	"net/http"

	"github.com/nvdaz/find-a-friend-api/db"
	"github.com/nvdaz/find-a-friend-api/service"

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
	// user := RegisterUserRequest{}
	// if err := c.Bind(&user); err != nil {
	// 	fmt.Println("Error parsing", err)
	// 	return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
	// }

	// if err := handler.userService.RegisterUser(user); err != nil {
	// 	fmt.Println("Error creating", err)
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "error creating user")
	// }

	return c.NoContent(http.StatusCreated)
}

func (handler *Handler) UpdateUser(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

type RegisterUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (handler *Handler) RegisterUser(c echo.Context) error {
	request := RegisterUserRequest{}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
	}

	user, err := handler.userService.RegisterUser(service.RegisterUser(request))
	if err != nil {
		fmt.Println("Error registering user", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error registering user")
	}

	return c.JSON(http.StatusOK, user)
}

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (handler *Handler) LoginUser(c echo.Context) error {
	request := LoginUserRequest{}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
	}

	user, err := handler.userService.LoginUser(service.LoginUser(request))
	if err != nil {
		fmt.Println("Error logging in user", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error logging in user")
	}

	return c.JSON(http.StatusOK, user)
}

type UpdateUserIconRequest struct {
	Id   string `json:"id"`
	Icon string `json:"icon"`
}

func (handler *Handler) UpdateUserIcon(c echo.Context) error {
	request := UpdateUserIconRequest{}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
	}

	err := handler.userService.UpdateAvatar(request.Id, request.Icon)
	if err != nil {
		fmt.Println("Error updating user icon", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error updating user icon")
	}

	return c.NoContent(http.StatusOK)

}
