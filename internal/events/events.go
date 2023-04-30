package events

import (
	"context"
	"github.com/eskpil/homehubsocial/pkg/database"
	"github.com/eskpil/homehubsocial/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var Queue []models.Event

func Get(c echo.Context) error {
	tmp := Queue
	Queue = []models.Event{}

	return c.JSON(http.StatusOK, tmp)
}

func Create(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	var body struct {
		Content string `json:"content"`
	}

	if err := c.Bind(&body); err != nil {
		log.Errorf("failed to bind body: %v", err)
		return err
	}

	event := models.Event{
		Id:      uuid.New().String(),
		Content: body.Content,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, err := database.Collection("events").InsertOne(ctx, &event); err != nil {
		log.Errorf("failed to insert new event: %v", err)
		return err
	}

	Queue = append(Queue, event)

	return c.JSON(http.StatusOK, event)
}
