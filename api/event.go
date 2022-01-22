package api

import (
	"github.com/asaberwd/event-messages/internal/event"
	"github.com/labstack/echo/v4"

	"net/http"
)

// EventHandler ...
type EventHandler struct {
	EventManager event.Manager
}

// NewEventHandler ...
func NewEventHandler(eventManager event.Manager) *EventHandler {
	return &EventHandler{EventManager: eventManager}
}

// Read ...
func (e *EventHandler) Read(ctx echo.Context) error {
	topic := ctx.Param("topic")
	res, err := e.EventManager.Read(topic)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, res)
}

// Write ...
func (e *EventHandler) Write(ctx echo.Context) error {
	topic := ctx.Param("topic")
	body := ""
	if err := ctx.Bind(&body); err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err := e.EventManager.Write(body, topic)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.String(http.StatusOK, "Created!")
}

// Router ...
func Router(e *echo.Echo, handler *EventHandler) {
	e.POST("/write/:topic", handler.Write)
	e.GET("/read/:topic", handler.Read)
}
