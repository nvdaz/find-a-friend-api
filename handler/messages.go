package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateMessageRequest struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Message    string `json:"message"`
}

func (handler *Handler) CreateMessage(c echo.Context) error {
	request := CreateMessageRequest{}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
	}

	err := handler.messageService.CreateMessage(request.SenderId, request.ReceiverId, request.Message)

	if err != nil {
		fmt.Println("Error creating message", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error creating message")
	}

	return c.NoContent(http.StatusCreated)
}

type GetMessagesRequest struct {
	UserId  string `json:"user_id"`
	OtherId string `json:"other_id"`
}

func (handler *Handler) GetMessages(c echo.Context) error {
	request := GetMessagesRequest{}
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
	}

	messages, err := handler.messageService.GetMessages(request.UserId, request.OtherId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error getting messages")
	}

	return c.JSON(http.StatusOK, messages)
}
