package handler

import (
	"net/http"

	"github.com/nvdaz/find-a-friend-api/db"

	"github.com/labstack/echo/v4"
)

func (handler *Handler) CreateServiceConversations(c echo.Context) error {
	serviceConversations := make([]db.ServiceConversation, 0)
	if err := c.Bind(&serviceConversations); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "error parsing request body")
	}

	if err := handler.serviceConversationStore.CreateServiceConversations(serviceConversations); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "error creating service conversations")
	}

	return c.NoContent(http.StatusCreated)
}

func (handler *Handler) GetServiceConversations(c echo.Context) error {
	id := c.Param("id")
	serviceConversations, err := handler.serviceConversationStore.GetRecentServiceConversations(id, 20)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, serviceConversations)
}
