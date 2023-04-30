package patients

import (
	"context"
	"fmt"
	"github.com/eskpil/homehubsocial/pkg/database"
	"github.com/eskpil/homehubsocial/pkg/models"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

func GetOne(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		log.Errorf("missing parameter id")
		return fmt.Errorf("missing param id")
	}

	var patient models.Patient
	if err := database.Collection("patients").FindOne(ctx, bson.D{{"_id", id}}).Decode(&patient); err != nil {
		log.Errorf("failed to find patient: %s: %v", id, err)
		return err
	}

	return c.JSON(http.StatusOK, patient)
}

func GetMany(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	cursor, err := database.Collection("patients").Find(ctx, bson.D{})
	if err != nil {
		return err
	}

	var patients []models.Patient
	for cursor.Next(ctx) {
		var patient models.Patient
		if err := cursor.Decode(&patient); err != nil {
			return err
		}

		patients = append(patients, patient)
	}

	return c.JSON(http.StatusOK, patients)
}
